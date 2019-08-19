package middlewares

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type (
	AsiraKafkaHandlers struct {
		KafkaConsumer     sarama.Consumer
		PartitionConsumer sarama.PartitionConsumer
	}

	Loan struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}
)

func init() {
	topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.kafka.topics", asira.App.ENV))

	kafka := &AsiraKafkaHandlers{}
	kafka.KafkaConsumer = asira.App.Kafka.Consumer

	kafka.SetPartitionConsumer(topics["loan_status_updt"].(string))

	go func() {
		for {
			message, err := kafka.Listen()
			if err != nil {
				log.Printf("error occured when listening kafka : %v", err)
			}
			if message != nil {
				err = loanUpdate(message)
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

func loanUpdate(kafkaMessage []byte) (err error) {
	var loanData Loan
	loan := models.Loan{}

	err = json.Unmarshal(kafkaMessage, &loanData)
	if err != nil {
		return err
	}

	data, err := loan.FindbyID(loanData.ID)
	if err != nil {
		return err
	}

	data.Status = loanData.Status
	_, err = data.Save()
	return err
}
