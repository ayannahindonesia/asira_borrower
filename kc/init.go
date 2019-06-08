package kc

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

type (
	App struct {
		Name    string        `json:"name"`
		Version string        `json:"version"`
		ENV     string        `json:"env"`
		Config  viper.Viper   `json:"prog_config"`
		DB      *gorm.DB      `json:"db"`
		Redis   *redis.Client `json:"redis"`
	}
)

// Initiate kayacredit instances
func init() {
	var err error
	Prog := &App{}
	Prog.Name = "kayacredit"
	Prog.Version = os.Getenv("APPVER")
	Prog.loadENV()
	if err = Prog.LoadConfigs(); err != nil {
		log.Println(err)
	}
	if err = Prog.DBinit(); err != nil {
		log.Println(err)
	}
	if err = Prog.RedisInit(); err != nil {
		log.Println(err)
	}
}

func (x *App) Close() error {
	var err error
	err = x.DB.Close()

	if err != nil {
		return err
	}
	return nil
}

// Loads environtment setting
func (x *App) loadENV() {
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
func (x *App) LoadConfigs() error {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	conf.AutomaticEnv()
	conf.SetConfigName("config")
	conf.AddConfigPath("$GOPATH/src/kayacredit")
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
func (x *App) DBinit() error {
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

	db.Exec(fmt.Sprintf("SET TIMEZONE TO '%s'", x.Config.GetString(fmt.Sprintf("%s.timezone", x.ENV))))
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(dbconf["maxlifetime"].(int)))
	db.DB().SetMaxIdleConns(dbconf["idle_conns"].(int))
	db.DB().SetMaxOpenConns(dbconf["open_conns"].(int))

	x.DB = db

	return nil
}

// Loads redis config
func (x *App) RedisInit() error {
	redisconf := x.Config.GetStringMap(fmt.Sprintf("%s.redis", x.ENV))
	db := redisconf["db"].(int)
	pass := redisconf["pass"].(string)
	host := redisconf["host"].(string)
	port := redisconf["port"].(string)

	address := host + ":" + port

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pass,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Failed Connect to Redis Server | " + err.Error())
	}

	x.Redis = client

	return nil
}
