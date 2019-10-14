package middlewares

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type (
	AsiraKafkaHandlers struct {
		KafkaConsumer     sarama.Consumer
		PartitionConsumer sarama.PartitionConsumer
	}

	Loan struct {
		ID           int       `json:"id"`
		Status       string    `json:"status"`
		DisburseDate time.Time `json:"disburse_date"`
	}
)

var wg sync.WaitGroup

func init() {
	var err error
	topics := asira.App.Config.GetStringMap(fmt.Sprintf("%s.kafka.topics.consumes", asira.App.ENV))

	kafka := &AsiraKafkaHandlers{}
	kafka.KafkaConsumer, err = sarama.NewConsumer([]string{asira.App.Kafka.Host}, asira.App.Kafka.Config)
	if err != nil {
		log.Printf("error while creating new kafka consumer : %v", err)
	}

	kafka.SetPartitionConsumer(topics["for_borrower"].(string))

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer kafka.KafkaConsumer.Close()
		for {
			message, err := kafka.Listen()
			if err != nil {
				log.Printf("error occured when listening kafka : %v", err)
			}
			if message != nil {
				err = processMessage(message)
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

func processMessage(kafkaMessage []byte) (err error) {
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
				err = bankType.FindbyID(ID)
				if err != nil {
					return err
				}

				err = bankType.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &bankType)
				if err != nil {
					return err
				}
				bankType.Save()
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
				err = bank.FindbyID(ID)
				if err != nil {
					return err
				}

				err = bank.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &bank)
				if err != nil {
					return err
				}
				bank.Save()
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
				err := bankService.FindbyID(ID)
				if err != nil {
					return err
				}

				err = bankService.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &bankService)
				if err != nil {
					return err
				}
				err = bankService.Save()
				return err
			}

		}
	case "bank_product":
		{
			var bankProduct models.BankProduct
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				err := bankProduct.FindbyID(ID)
				if err != nil {
					return err
				}

				err = bankProduct.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &bankProduct)
				if err != nil {
					return err
				}
				err = bankProduct.Save()
				return err
			}

		}
	case "loan":
		log.Printf("message : %v", string(kafkaMessage))
		err = loanUpdate([]byte(data[1]))
		if err != nil {
			return err
		}
	case "loan_purpose":
		{
			var LoanPurpose models.LoanPurpose
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				err := LoanPurpose.FindbyID(ID)
				if err != nil {
					return err
				}

				err = LoanPurpose.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &LoanPurpose)
				if err != nil {
					return err
				}
				err = LoanPurpose.Save()
				return err
			}

		}
	default:
		return nil
	}
	return nil
}

func loanUpdate(kafkaMessage []byte) (err error) {
	var loanData Loan
	loan := models.Loan{}

	err = json.Unmarshal(kafkaMessage, &loanData)
	if err != nil {
		return err
	}

	err = loan.FindbyID(loanData.ID)
	if err != nil {
		return err
	}

	loan.Status = loanData.Status
	loan.DisburseDate = loanData.DisburseDate
	err = loan.Save()
	return err
}
