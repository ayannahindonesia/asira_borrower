package asira

import (
	"asira_borrower/cron"
	"asira_borrower/custommodule"
	"asira_borrower/validator"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/ayannahindonesia/basemodel"
	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"github.com/xlzd/gotp"

	"github.com/ayannahindonesia/northstar/lib/northstarlib"
)

var (
	App *Application
)

type (
	Application struct {
		Name      string        `json:"name"`
		Port      string        `json:"port"`
		Version   string        `json:"version"`
		ENV       string        `json:"env"`
		Config    viper.Viper   `json:"prog_config"`
		DB        *gorm.DB      `json:"db"`
		OTP       OTP           `json:"otp"`
		Kafka     KafkaInstance `json:"kafka"`
		Messaging custommodule.Messaging
		S3        custommodule.S3      `json:"s3"`
		Emailer   custommodule.Emailer `json:"email"`
		Northstar northstarlib.NorthstarLib
		Cron      cron.Cron `json:"cron"`
	}

	OTP struct {
		HOTP *gotp.HOTP
		TOTP *gotp.TOTP
	}

	KafkaInstance struct {
		Config *sarama.Config
		Host   string
	}
)

// Initiate asira instances
func init() {
	var err error
	App = &Application{}
	App.Name = "asira_borrower"
	App.Port = os.Getenv("APPPORT")
	App.Version = os.Getenv("APPVER")
	App.loadENV()
	if err = App.LoadConfigs(); err != nil {
		log.Printf("Load config error : %v", err)
	}
	if err = App.DBinit(); err != nil {
		log.Printf("DB init error : %v", err)
	}

	App.KafkaInit()
	App.MessagingInit()
	App.S3init()
	App.EmailerInit()
	App.NorthstarInit()
	App.CronInit()

	otpSecret := gotp.RandomSecret(16)
	App.OTP = OTP{
		HOTP: gotp.NewDefaultHOTP(otpSecret),
		TOTP: gotp.NewDefaultTOTP(otpSecret),
	}

	// apply custom validator
	v := validator.AsiraValidator{DB: App.DB}
	v.CustomValidatorRules()
}

func (x *Application) Close() (err error) {
	if err = x.DB.Close(); err != nil {
		return err
	}

	return nil
}

// Loads environtment setting
func (x *Application) loadENV() {
	APPENV := os.Getenv("APPENV")

	switch APPENV {
	default:
		x.ENV = "development"
		break
	case "development":
		x.ENV = "development"
		break
	case "staging":
		x.ENV = "staging"
		break
	case "production":
		x.ENV = "production"
		break
	}
}

// Loads general configs
func (x *Application) LoadConfigs() error {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	conf.AutomaticEnv()
	conf.SetConfigName("config")
	conf.AddConfigPath(os.Getenv("CONFIGPATH"))
	conf.SetConfigType("yaml")
	if err := conf.ReadInConfig(); err != nil {
		return err
	}
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Println("App Config file changed %s:", e.Name)
		x.LoadConfigs()
	})
	x.Config = viper.Viper(*conf)

	return nil
}

// Loads DBinit configs
func (x *Application) DBinit() error {
	dbconf := x.Config.GetStringMap(fmt.Sprintf("%s.database", x.ENV))
	Cons := basemodel.DBConfig{
		Adapter:        basemodel.PostgresAdapter,
		Host:           dbconf["host"].(string),
		Port:           dbconf["port"].(string),
		Username:       dbconf["username"].(string),
		Password:       dbconf["password"].(string),
		Table:          dbconf["table"].(string),
		Timezone:       dbconf["timezone"].(string),
		Maxlifetime:    dbconf["maxlifetime"].(int),
		IdleConnection: dbconf["idle_conns"].(int),
		OpenConnection: dbconf["open_conns"].(int),
		SSL:            dbconf["sslmode"].(string),
		Logmode:        dbconf["logmode"].(bool),
	}
	basemodel.Start(Cons)
	x.DB = basemodel.DB
	return nil
}

func (x *Application) KafkaInit() {
	kafkaConf := x.Config.GetStringMap(fmt.Sprintf("%s.kafka", x.ENV))

	if kafkaConf["log_verbose"].(bool) {
		sarama.Logger = log.New(os.Stdout, "[borrower kafka] ", log.LstdFlags)
	}

	x.Kafka.Config = sarama.NewConfig()
	x.Kafka.Config.ClientID = kafkaConf["client_id"].(string)
	if kafkaConf["sasl"].(bool) {
		x.Kafka.Config.Net.SASL.Enable = true
	}
	x.Kafka.Config.Net.SASL.User = kafkaConf["user"].(string)
	x.Kafka.Config.Net.SASL.Password = kafkaConf["pass"].(string)
	x.Kafka.Config.Producer.Return.Successes = true
	x.Kafka.Config.Producer.Partitioner = sarama.NewRandomPartitioner
	x.Kafka.Config.Producer.RequiredAcks = sarama.WaitForAll
	x.Kafka.Config.Consumer.Return.Errors = true
	x.Kafka.Host = strings.Join([]string{kafkaConf["host"].(string), kafkaConf["port"].(string)}, ":")
}

// MessagingInit func
func (x *Application) MessagingInit() {
	messagingConfig := x.Config.GetStringMap(fmt.Sprintf("%s.messaging", x.ENV))

	endpoints := custommodule.MessagingEndpoints{
		ClientAuth:       messagingConfig["client_auth"].(string),
		SMS:              messagingConfig["sms"].(string),
		PushNotification: messagingConfig["push_notification"].(string),
		ListNotification: messagingConfig["list_notification"].(string),
	}

	x.Messaging.SetConfig(messagingConfig["key"].(string), messagingConfig["secret"].(string), messagingConfig["url"].(string), endpoints)
}

// S3init load config for s3
func (x *Application) S3init() (err error) {
	s3conf := x.Config.GetStringMap(fmt.Sprintf("%s.s3", x.ENV))

	x.S3, err = custommodule.NewS3(s3conf["access_key"].(string), s3conf["secret_key"].(string), s3conf["host"].(string), s3conf["bucket_name"].(string), s3conf["region"].(string))

	return err
}

//EmailerInit load config for s3
func (x *Application) EmailerInit() (err error) {
	emailerConf := x.Config.GetStringMap(fmt.Sprintf("%s.mailer", x.ENV))

	x.Emailer = custommodule.Emailer{
		Host:     emailerConf["host"].(string),
		Port:     emailerConf["port"].(int),
		Email:    emailerConf["email"].(string),
		Password: emailerConf["password"].(string),
	}
	return err
}

// NorthstarInit config for northstar logger
func (x *Application) NorthstarInit() {
	northstarconf := x.Config.GetStringMap(fmt.Sprintf("%s.northstar", x.ENV))

	x.Northstar = northstarlib.NorthstarLib{
		Host:         App.Kafka.Host,
		Secret:       northstarconf["secret"].(string),
		Topic:        northstarconf["topic"].(string),
		Send:         northstarconf["send"].(bool),
		SaramaConfig: App.Kafka.Config,
	}
}

// CronInit load cron
func (x *Application) CronInit() (err error) {
	x.Cron.TZ = x.Config.GetString(fmt.Sprintf("%s.database.timezone", x.ENV))
	cron.DB = x.DB
	x.Cron.Time = x.Config.GetString(fmt.Sprintf("%s.cron.time", x.ENV))
	x.Cron.New()
	x.Cron.Start()

	return nil
}
