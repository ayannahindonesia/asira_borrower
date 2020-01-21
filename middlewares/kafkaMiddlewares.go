package middlewares

import (
	"asira_borrower/asira"
	"asira_borrower/handlers"
	"asira_borrower/models"
	"errors"

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
		ID                  uint64    `json:"id"`
		Status              string    `json:"status"`
		DisburseDate        time.Time `json:"disburse_date"`
		DisburseDateChanged bool      `json:"disburse_date_changed"`
		DisburseStatus      string    `json:"disburse_status"`
		RejectReason        string    `json:"reject_reason"`
		DueDate             time.Time `json:"due_date"`
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

func handleOperation(modObj interface{}, mode interface{}) error{
	var err error
	//get object
	switch modObj.(type) {
	case models.BankType:
		mod := modObj.(models.BankType)
		break
	case models.Bank:
		mod := modObj.(models.Bank)
		break
	case models.Service:
		mod := modObj.(models.Service)
		break
	case models.Product:
		mod := modObj.(models.Product)
		break
	case models.LoanPurpose:
		mod := modObj.(models.LoanPurpose)
		break
	}

	//
	switch mode.(string) {
	default:
		err = fmt.Errorf("invalid payload")
		break
	case "create":
		err = mod.Create()
		break
	case "update":
		err = mod.Save()
		break
	case "delete":
		err = mod.Delete()
		break
	}
	
	return err
}

func processMessage(kafkaMessage []byte) (err error) {

	var arr map[string]interface{}

	//parse msg
	data := strings.SplitN(string(kafkaMessage), ":", 2)
	err = json.Unmarshal([]byte(data[1]), &arr)
	if err != nil {
		return err
	}
	marshal, _ := json.Marshal(arr["payload"])

	//cek obj type
	switch data[0] {
	case "bank_type":
		{
			mod := models.BankType{}
			json.Unmarshal(marshal, &mod)
			handleOperation(mod, arr["mode"])
		}
		break
	case "bank":
		{
			mod := models.Bank{}
			json.Unmarshal(marshal, &mod)
			handleOperation(mod, arr["mode"])
		}
		break
	case "service":
		{
			mod := models.Service{}
			json.Unmarshal(marshal, &mod)
			handleOperation(mod, arr["mode"])
		}
		break
	case "product":
		{	
			mod := models.Product{}
			json.Unmarshal(marshal, &mod)
			handleOperation(mod, arr["mode"])
		}
		break
	case "loan_purpose":
		{
			mod := models.LoanPurpose{}
			json.Unmarshal(marshal, &mod)
			handleOperation(mod, arr["mode"])
		}
		break
	case "loan":
		// log.Printf("message : %v", string(kafkaMessage))
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
	type Filter struct {
		ID                  uint64    `json:"id"`
		Status              string    `json:"status"`
		DisburseDate        time.Time `json:"disburse_date"`
		DisburseStatus      string    `json:"disburse_status"`
		DisburseDateChanged bool      `json:"disburse_date_changed"`
	}
	var loanData Loan
	loan := models.Loan{}
	borrower := models.Borrower{}

	err = json.Unmarshal(kafkaMessage, &loanData)
	if err != nil {
		return err
	}

	err = loan.FilterSearchSingle(&Filter{
		ID:                  loanData.ID,
		Status:              loanData.Status,
		DisburseDate:        loanData.DisburseDate,
		DisburseStatus:      loanData.DisburseStatus,
		DisburseDateChanged: loanData.DisburseDateChanged,
	})
	//data ada di kafka sebelumnya
	if err == nil {
		return errors.New("loan already in db")
	}

	//get by ID saja
	err = loan.FindbyID(loanData.ID)
	if err != nil {
		return err
	}

	//copy data
	loan.Status = loanData.Status
	loan.DisburseDate = loanData.DisburseDate
	loan.DisburseStatus = loanData.DisburseStatus
	loan.DisburseDateChanged = loanData.DisburseDateChanged
	loan.RejectReason = loanData.RejectReason
	loan.DueDate = loanData.DueDate
	err = loan.SaveNoKafka()
	if err != nil {
		return err
	}

	err = borrower.FindbyID(loan.Borrower)
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

	to := borrower.Email
	subject := "[NO REPLY] - " + title
	link := "" //FUTURE: link open apps detail
	message := formatedMsg + link

	err = handlers.SendMail(to, subject, message)
	if err != nil {
		log.Println(err.Error())
	}

	//get user login n fcm data from borrower
	user := models.User{}
	err = user.FindbyBorrowerID(borrower.ID)
	if err != nil {
		log.Println(err.Error())
	}

	//send notif
	fmt.Println("FCMToken : ", user.FCMToken)
	responseBody, err := asira.App.Messaging.SendNotificationByToken(title, formatedMsg, mapData, user.FCMToken, recipientID)
	if err != nil {
		type ErrorResponse struct {
			Details string `json:"details"`
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse

		//parse error response
		err = json.Unmarshal(responseBody, &errorResponse)
		if err != nil {
			log.Printf(err.Error())
			return err
		}

		//set error notif
		notif := models.Notification{}
		notif.Title = "failed"
		notif.ClientID = 2
		notif.RecipientID = recipientID
		notif.Response = errorResponse.Message
		err = notif.Create()
		return err
	}

	log.Println("Response Body : ", string(responseBody))
	//logging notification
	var notif models.Notification
	err = json.Unmarshal(responseBody, &notif)
	if err != nil {
		msg := fmt.Sprintf("%s : %s", "Failed storing notification", err.Error())
		log.Printf(msg)
		return fmt.Errorf(msg)
	} else {
		notif.Create()
	}

	return err
}

type Filter struct {
	Username string `json:"username"`
}

func syncAgent(dataAgent []byte) (err error) {

	var agent models.Agent
	var a map[string]interface{}
	fmt.Println("dataAgent => ", dataAgent)
	err = json.Unmarshal(dataAgent, &a)
	if err != nil {
		return err
	}

	if a["delete"] != nil && a["delete"].(bool) == true {
		ID := uint64(a["id"].(float64))
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
		err = agent.SaveNoKafka()

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
		borrower models.Borrower
		bank     models.Bank
	)
	//get loan
	Loan := object.(models.Loan)

	//get bank
	borrower.FindbyID(Loan.Borrower)
	bank.FindbyID(uint64(borrower.Bank.Int64))

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
