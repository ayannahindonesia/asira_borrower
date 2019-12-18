package handlers

import (
	"asira_borrower/asira"
	"fmt"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// Info main type
type Info struct {
	Time   string `json:"time"`
	Stacks bool   `json:"stacks"`
}

var (
	err  error
	info Info
)

// ServiceInfo check service info
func ServiceInfo(c echo.Context) error {
	defer c.Request().Body.Close()

	info.Time = fmt.Sprintf("%v", time.Now().Format("2006-01-02T15:04:05"))
	info.Stacks = true
	if err = healthcheckKafka(); err != nil {
		info.Stacks = false
	}
	if err = healthcheckDB(); err != nil {
		info.Stacks = false
	}

	return c.JSON(http.StatusOK, info)
}
func healthcheckKafka() (err error) {
	producer, err := sarama.NewAsyncProducer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		return err
	}
	defer producer.Close()

	consumer, err := sarama.NewConsumer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		return err
	}
	defer consumer.Close()

	return nil
}

func healthcheckDB() (err error) {
	dbconf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.database", asira.App.ENV))
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=5", dbconf["username"].(string), dbconf["password"].(string), dbconf["host"].(string), dbconf["port"].(string), dbconf["table"].(string), dbconf["sslmode"].(string))

	db, err := gorm.Open("postgres", connectionString)
	defer db.Close()
	return err
}
