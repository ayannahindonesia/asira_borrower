package middlewares

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
)

type (
	AsiraKafkaHandlers struct {
		KafkaConsumer     sarama.Consumer
		PartitionConsumer sarama.PartitionConsumer
	}
)

func init() {
	topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.kafka.topics", asira.App.ENV))

	kafka := &AsiraKafkaHandlers{}
	kafka.KafkaConsumer = asira.App.Kafka.Consumer

	kafka.SetPartitionConsumer(topics["entity_hook"].(string))

	go func() {
		for {
			message, err := kafka.Listen()
			if err != nil {
				log.Printf("error occured when listening kafka : %v", err)
			}
			if message != nil {
				err = getEntity(message)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

func (k *AsiraKafkaHandlers) SetPartitionConsumer(topic string) (err error) {
	k.PartitionConsumer, err = k.KafkaConsumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}

	return nil
}

func (k *AsiraKafkaHandlers) Listen() ([]byte, error) {
	select {
	case err := <-k.PartitionConsumer.Errors():
		return nil, err
	case msg := <-k.PartitionConsumer.Messages():
		return msg.Value, nil
	}

	return nil, fmt.Errorf("unidentified error while listening")
}

func getEntity(kafkaMessage []byte) (err error) {
	data := strings.SplitN(string(kafkaMessage), ":", 2)
	switch data[0] {
	case "bank_type":
		{
			var bankTData struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}

			var bankType models.BankType
			err = json.Unmarshal([]byte(data[1]), &bankTData)
			if err != nil {
				return err
			}

			data, err := bankType.FindbyID(bankTData.ID)
			if err != nil {
				return err
			}
			data.Name = bankTData.Name
			_, err = data.Save()
			return err
		}
	default:
		{
			log.Println(data[1])
		}
	}
	return err
}
