package middlewares

import (
	"asira_borrower/asira"
	"asira_borrower/handlers"
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
		ID                  int       `json:"id"`
		Status              string    `json:"status"`
		DisburseDate        time.Time `json:"disburse_date"`
		DisburseDateChanged bool      `json:"disburse_date_changed"`
		DisburseStatus      string    `json:"disburse_status"`
		RejectReason        string    `json:"reject_reason"`
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
		break
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
		break
	case "service":
		{
			var service models.Service
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				err := service.FindbyID(ID)
				if err != nil {
					return err
				}

				err = service.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &service)
				if err != nil {
					return err
				}
				err = service.Save()
				return err
			}
		}
		break
	case "product":
		{
			var product models.Product
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				err := product.FindbyID(ID)
				if err != nil {
					return err
				}

				err = product.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &product)
				if err != nil {
					return err
				}
				err = product.Save()
				return err
			}
		}
		break
	case "loan_purpose":
		{
			var loanPurpose models.LoanPurpose
			var a map[string]interface{}

			err = json.Unmarshal([]byte(data[1]), &a)
			if err != nil {
				return err
			}

			if a["delete"] != nil && a["delete"].(bool) == true {
				ID := int(a["id"].(float64))
				err := loanPurpose.FindbyID(ID)
				if err != nil {
					return err
				}

				err = loanPurpose.Delete()
				if err != nil {
					return err
				}
			} else {
				err = json.Unmarshal([]byte(data[1]), &loanPurpose)
				if err != nil {
					return err
				}
				err = loanPurpose.Save()
				return err
			}
		}
		break
	case "loan":
		log.Printf("message : %v", string(kafkaMessage))
		err = loanUpdate([]byte(data[1]))
		if err != nil {
			return err
		}
		break
	case "agent":
		log.Printf("message : %v", string(kafkaMessage))
		err = syncAgent([]byte(data[1]))
		if err != nil {
			return err
		}
		break
	default:
		return nil
		break
	}
	return nil
}

func loanUpdate(kafkaMessage []byte) (err error) {
	var loanData Loan
	loan := models.Loan{}
	borrower := models.Borrower{}

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
	loan.DisburseStatus = loanData.DisburseStatus
	loan.DisburseDateChanged = loanData.DisburseDateChanged
	loan.RejectReason = loanData.RejectReason
	err = loan.SaveNoKafka()

	err = borrower.FindbyID(int(loan.Owner.Int64))
	if err != nil {
		return err
	}
	fmt.Printf("loanData =>> %+v", loanData)
	var formatedMsg string
	if loanData.DisburseStatus == "processing" && loanData.DisburseDateChanged == true {
		//tgl pencairan diubah oleh pihak bank
		formatedMsg = FormatingMessage("disburse_changed", loan)

	} else if loanData.Status == "approved" && loanData.DisburseStatus == "processing" {
		//pinjaman diterima oleh bank
		formatedMsg = FormatingMessage("loan_approved", loan)

	} else if loanData.Status == "approved" && loanData.DisburseStatus == "confirmed" {
		//pinjaman telah dicairkan
		formatedMsg = FormatingMessage("disburse", loan)

	} else if loanData.Status == "rejected" {
		//pinjaman ditolak oleh bank
		formatedMsg = FormatingMessage("loan_rejected", loan)
	}

	//custom map data for firebase key "Data"
	mapData := map[string]string{
		"id":     fmt.Sprintf("%d", loan.ID),
		"status": loan.Status,
	}

	//set recipient ID
	recipientID := fmt.Sprintf("borrower-%d", borrower.ID)

	//set title for notif and email
	title := "Status Pinjaman Anda"

	//send notif
	err = asira.App.Messaging.SendNotificationByToken(title, formatedMsg, mapData, borrower.FCMToken, recipientID)
	if err != nil {
		return err
	}

	to := borrower.Email
	subject := "[NO REPLY] - " + title
	link := "" //FUTURE: link open apps detail
	message := formatedMsg + link

	err = handlers.SendMail(to, subject, message)
	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func syncAgent(dataAgent []byte) (err error) {
	var agent models.Agent
	var a map[string]interface{}

	err = json.Unmarshal(dataAgent, &a)
	if err != nil {
		return err
	}

	if a["delete"] != nil && a["delete"].(bool) == true {
		ID := int(a["id"].(float64))
		err := agent.FindbyID(ID)
		if err != nil {
			return err
		}

		err = agent.Delete()
		if err != nil {
			return err
		}
	} else {
		err = json.Unmarshal(dataAgent, &agent)
		if err != nil {
			return err
		}
		err = agent.Save()
		return err
	}
	return nil
}
func FormatingMessage(msgType string, object interface{}) string {

	var msg string

	var (
		status string
		prefix string
		//postfix string
		owner models.Borrower
		bank  models.Bank
	)
	//get loan
	Loan := object.(models.Loan)

	//get bank
	owner.FindbyID(int(Loan.Owner.Int64))
	bank.FindbyID(int(owner.Bank.Int64))

	//NOTE format pesan (PRD 7)
	// format := "Loan id %d %s oleh %s. "
	format := "Pinjaman nomor %d %s oleh %s, silahkan cek di aplikasi."
	// approvedFormat := "Dapat dicairkan pada %s"

	switch msgType {
	case "loan_approved":
		status = "diterima"
		format = prefix + format                              // + postfix
		msg = fmt.Sprintf(format, Loan.ID, status, bank.Name) //, Loan.DisburseDate)
		break

	case "loan_rejected":
		prefix = "Maaf, "
		status = "ditolak"
		format = prefix + format                              // + postfix
		msg = fmt.Sprintf(format, Loan.ID, status, bank.Name) //, Loan.DisburseDate)
		break

	case "disburse":
		format := "Pinjaman nomor %d dari bank %s telah dicairkan, silahkan cek di aplikasi."
		msg = fmt.Sprintf(format, Loan.ID, bank.Name) //, Loan.DisburseDate)
		break

	case "disburse_changed":
		format := "Maaf, tanggal pencairan untuk pinjaman %d dari bank %s direvisi menjadi tanggal %s."
		msg = fmt.Sprintf(format, Loan.ID, bank.Name, Loan.DisburseDate)
		break

	}

	return msg
}
