package handlers

import (
	"asira_borrower/asira"
	"fmt"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var err error

func Ping(c echo.Context) error {
	defer c.Request().Body.Close()

	if err = healthcheckKafka(); err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, nil, "Server is not ready")
	}
	if err = healthcheckDB(); err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, nil, "Server is not ready")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Server ready"})
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
