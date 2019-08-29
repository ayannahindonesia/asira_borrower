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

	Loan struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}
)

func init() {
	topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.kafka.topics.consumes", asira.App.ENV))

	kafka := &AsiraKafkaHandlers{}
	kafka.KafkaConsumer = asira.App.Kafka.Consumer

	kafka_loan := &AsiraKafkaHandlers{}
	kafka_loan.KafkaConsumer = asira.App.Kafka.Consumer

	kafka.SetPartitionConsumer(topics["from_lender"].(string))
	kafka_loan.SetPartitionConsumer(topics["loan_status_updt"].(string))

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

			message_loan, err := kafka_loan.Listen()
			if err != nil {
				log.Printf("error occured when listening kafka : %v", err)
			}
			if message_loan != nil {
				err = loanUpdate(message_loan)
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
			var bankType models.BankType
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				result, err := bankType.FindbyID(ID)
				if err != nil {
					return err
				}

				_, err = result.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &bankType)
				if err != nil {
					return err
				}
				_, err = bankType.Save()
				return err
			}

		}

	case "bank":
		{
			var bank models.Bank
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				result, err := bank.FindbyID(ID)
				if err != nil {
					return err
				}

				_, err = result.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &bank)
				if err != nil {
					return err
				}
				_, err = bank.Save()
				return err
			}

		}
	case "bank_service":
		{
			var bankService models.BankService
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				result, err := bankService.FindbyID(ID)
				if err != nil {
					return err
				}

				_, err = result.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &bankService)
				if err != nil {
					return err
				}
				_, err = bankService.Save()
				return err
			}

		}
	case "service_product":
		{
			var serviceProduct models.ServiceProduct
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				result, err := serviceProduct.FindbyID(ID)
				if err != nil {
					return err
				}

				_, err = result.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &serviceProduct)
				if err != nil {
					return err
				}
				_, err = serviceProduct.Save()
				return err
			}

		}
	default:
		{
			return err
		}
	}
	return err
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
