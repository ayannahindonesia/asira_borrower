package asira

import (
	"asira_borrower/validator"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"github.com/xlzd/gotp"
)

var (
	App *Application
)

type (
	Application struct {
		Name    string        `json:"name"`
		Version string        `json:"version"`
		ENV     string        `json:"env"`
		Config  viper.Viper   `json:"prog_config"`
		DB      *gorm.DB      `json:"db"`
		OTP     OTP           `json:"otp"`
		Kafka   KafkaInstance `json:"kafka"`
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
	App.Version = os.Getenv("APPVER")
	App.loadENV()
	if err = App.LoadConfigs(); err != nil {
		log.Printf("Load config error : %v", err)
	}
	if err = App.DBinit(); err != nil {
		log.Printf("DB init error : %v", err)
	}

	App.KafkaInit()

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
	conf.AddConfigPath("$GOPATH/src/asira_borrower")
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

// Loads DB postgres configs
func (x *Application) DBinit() error {
	dbconf := x.Config.GetStringMap(fmt.Sprintf("%s.database", x.ENV))
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dbconf["username"].(string), dbconf["password"].(string), dbconf["host"].(string), dbconf["port"].(string), dbconf["table"].(string), dbconf["sslmode"].(string))

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	if err = db.DB().Ping(); err != nil {
		return err
	}

	db.LogMode(dbconf["logmode"].(bool))

	db.Exec(fmt.Sprintf("SET TIMEZONE TO '%s'", dbconf["timezone"].(string)))
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(dbconf["maxlifetime"].(int)))
	db.DB().SetMaxIdleConns(dbconf["idle_conns"].(int))
	db.DB().SetMaxOpenConns(dbconf["open_conns"].(int))

	x.DB = db

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
