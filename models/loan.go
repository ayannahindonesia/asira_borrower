package models

import (
	"asira_borrower/custommodule/irate"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Shopify/sarama"
	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

type (
	// Loan main struct
	Loan struct {
		basemodel.BaseModel
		Borrower            uint64         `json:"borrower" gorm:"column:borrower;foreignkey"`
		Status              string         `json:"status" gorm:"column:status;type:varchar(255)" sql:"DEFAULT:'processing'"`
		LoanAmount          float64        `json:"loan_amount" gorm:"column:loan_amount;type:int;not null"`
		Installment         int            `json:"installment" gorm:"column:installment;type:int;not null"` // plan of how long loan to be paid
		InstallmentDetails  pq.Int64Array  `json:"installment_details" gorm:"column:installment_details"`
		Fees                postgres.Jsonb `json:"fees" gorm:"column:fees;type:jsonb"`
		Interest            float64        `json:"interest" gorm:"column:interest;type:int;not null"`
		TotalLoan           float64        `json:"total_loan" gorm:"column:total_loan;type:int;not null"`
		DueDate             time.Time      `json:"due_date" gorm:"column:due_date"`
		DisburseAmount      float64        `json:"disburse_amount" gorm:"column:disburse_amount;type:int;not null"`
		LayawayPlan         float64        `json:"layaway_plan" gorm:"column:layaway_plan;type:int;not null"` // how much borrower will pay per month
		Product             uint64         `json:"product" gorm:"column:product;foreignkey"`
		LoanIntention       string         `json:"loan_intention" gorm:"column:loan_intention;type:varchar(255);not null"`
		IntentionDetails    string         `json:"intention_details" gorm:"column:intention_details;type:text;not null"`
		BorrowerInfo        postgres.Jsonb `json:"borrower_info" gorm:"column:borrower_info;type:jsonb"`
		OTPverified         bool           `json:"otp_verified" gorm:"column:otp_verified;type:boolean" sql:"DEFAULT:FALSE"`
		DisburseDate        time.Time      `json:"disburse_date" gorm:"column:disburse_date"`
		DisburseDateChanged bool           `json:"disburse_date_changed" gorm:"column:disburse_date_changed"`
		DisburseStatus      string         `json:"disburse_status" gorm:"column:disburse_status" sql:"DEFAULT:'processing'"`
		ApprovalDate        time.Time      `json:"approval_date" gorm:"column:approval_date"`
		RejectReason        string         `json:"reject_reason" gorm:"column:reject_reason"`
		FormInfo            postgres.Jsonb `json:"form_info" gorm:"column:form_info;type:jsonb"`
	}

	// LoanFee struct
	LoanFee struct {
		Description string `json:"description"`
		Amount      string `json:"amount"`
		FeeMethod   string `json:"fee_method"`
	}
	// LoanFees slice of LoanFee
	LoanFees []LoanFee
)

// BeforeCreate gorm callback hook
func (l *Loan) BeforeCreate() (err error) {
	borrower := Borrower{}
	err = borrower.FindbyID(l.Borrower)
	if err != nil {
		return err
	}

	jsonB, err := json.Marshal(borrower)
	if err != nil {
		return err
	}
	l.BorrowerInfo = postgres.Jsonb{jsonB}

	err = l.SetProductLoanReferences()
	if err != nil {
		return err
	}

	err = l.Calculate()
	if err != nil {
		log.Printf("calculate error : %+v", err)
		return err
	}

	return nil
}

// SetProductLoanReferences func
func (l *Loan) SetProductLoanReferences() (err error) {
	product := Product{}
	err = product.FindbyID(l.Product)
	if err != nil {
		return err
	}

	l.Fees = product.Fees
	l.Interest = product.Interest

	return nil
}

// Calculate func
func (l *Loan) Calculate() (err error) {
	// calculate total loan
	var (
		fee        float64
		fees       LoanFees
		borrower   Borrower
		bank       Bank
		parsedFees LoanFees
	)

	borrower.FindbyID(l.Borrower)
	bank.FindbyID(uint64(borrower.Bank.Int64))

	err = l.CalculateInterest()
	if err != nil {
		return err
	}
	l.DisburseAmount = l.LoanAmount

	json.Unmarshal(l.Fees.RawMessage, &fees)

	for _, v := range fees {
		if strings.ContainsAny(v.Amount, "%") {
			feeString := strings.TrimFunc(v.Amount, func(r rune) bool {
				return !unicode.IsNumber(r)
			})
			f, _ := strconv.Atoi(feeString)
			fee = (float64(f) / 100) * l.LoanAmount
		} else {
			f, _ := strconv.Atoi(v.Amount)
			fee = float64(f)
		}

		// parse fees
		parsedFees = append(parsedFees, LoanFee{
			Description: v.Description,
			Amount:      fmt.Sprintf("%f", fee),
			FeeMethod:   v.FeeMethod,
		})

		switch v.FeeMethod {
		case "deduct_loan":
			l.DisburseAmount -= fee
			break
		case "charge_loan":
			l.TotalLoan += fee
			l.LayawayPlan += fee / float64(l.Installment)
			break
		}
	}
	// parse fees
	jMarshal, _ := json.Marshal(parsedFees)
	l.Fees = postgres.Jsonb{jMarshal}

	return nil
}

// CalculateInterest func
func (l *Loan) CalculateInterest() (err error) {
	var product Product

	err = product.FindbyID(l.Product)
	switch product.InterestType {
	default:
		return err
	case "flat":
		err = l.FlatFormula(product.RecordInstallmentDetails)
		break
	case "onetimepay":
		err = l.OnetimepayFormula(product.RecordInstallmentDetails)
		break
	case "fixed":
		err = l.FixedInterestFormula(product.RecordInstallmentDetails)
		break
	case "efektif_menurun":
		err = l.EfektifMenurunFormula(product.RecordInstallmentDetails)
		break
	}

	return err
}

// FlatFormula func
func (l *Loan) FlatFormula(x bool) (err error) {
	var (
		pokok          float64
		bunga          float64
		installments   []Installment
		installmentsID []int64
	)

	pokok, bunga, l.LayawayPlan, l.TotalLoan = irate.FLATANNUAL(l.Interest/100, l.LoanAmount, float64(l.Installment))
	if x {
		for i := 1; i <= l.Installment; i++ {
			duedate := time.Now().AddDate(0, i, 0)
			installment := Installment{
				Period:          i,
				LoanPayment:     pokok,
				InterestPayment: bunga,
				DueDate:         &duedate,
			}
			err := installment.Create()
			if err != nil {
				return err
			}
			installments = append(installments, installment)
			installmentsID = append(installmentsID, int64(installment.ID))
		}
		err = syncInstallment(installments)
		l.InstallmentDetails = pq.Int64Array(installmentsID)
	}

	return err
}

// OnetimepayFormula func
func (l *Loan) OnetimepayFormula(x bool) (err error) {
	var (
		pokok          float64
		bunga          float64
		installments   []Installment
		installmentsID []int64
	)

	pokok, bunga, l.LayawayPlan, l.TotalLoan = irate.ONETIMEPAYMENT(l.Interest/100, l.LoanAmount, float64(l.Installment))
	if x {
		for i := 1; i <= l.Installment; i++ {
			duedate := time.Now().AddDate(0, i, 0)
			installment := Installment{
				Period:          i,
				LoanPayment:     pokok,
				InterestPayment: bunga,
				DueDate:         &duedate,
			}
			err := installment.Create()
			if err != nil {
				return err
			}
			installments = append(installments, installment)
			installmentsID = append(installmentsID, int64(installment.ID))
		}
		err = syncInstallment(installments)
		l.InstallmentDetails = pq.Int64Array(installmentsID)
	}
	return err
}

// FixedInterestFormula func
func (l *Loan) FixedInterestFormula(x bool) (err error) {
	rate := ((l.Interest / 100) / 12)
	var (
		pokok          float64
		bunga          float64
		installments   []Installment
		installmentsID []int64
	)
	for i := 1; i <= l.Installment; i++ {
		pokok, bunga = irate.PIPMT(rate, float64(i), float64(l.Installment), -l.LoanAmount, 1)
		duedate := time.Now().AddDate(0, i, 0)

		installment := Installment{
			Period:          i,
			LoanPayment:     pokok,
			InterestPayment: bunga,
			DueDate:         &duedate,
		}
		err := installment.Create()
		if err != nil {
			return err
		}
		installments = append(installments, installment)
		installmentsID = append(installmentsID, int64(installment.ID))
	}

	err = syncInstallment(installments)

	l.InstallmentDetails = pq.Int64Array(installmentsID)
	l.LayawayPlan = pokok + bunga
	l.TotalLoan = l.LayawayPlan * float64(l.Installment)
	return err
}

// EfektifMenurunFormula func
func (l *Loan) EfektifMenurunFormula(x bool) (err error) {
	plafon := l.LoanAmount
	cicilanpokok := l.LoanAmount / float64(l.Installment)
	var (
		cicilanbungas  []float64
		installments   []Installment
		installmentsID []int64
	)
	for i := 1; i <= l.Installment; i++ {
		bunga := plafon * (l.Interest / 100) / 12
		cicilanbungas = append(cicilanbungas, bunga)
		plafon -= cicilanpokok
		duedate := time.Now().AddDate(0, i, 0)

		installment := Installment{
			Period:          i,
			LoanPayment:     cicilanpokok,
			InterestPayment: bunga,
			DueDate:         &duedate,
		}
		err = installment.Create()
		if err != nil {
			return err
		}
		installments = append(installments, installment)
		installmentsID = append(installmentsID, int64(installment.ID))
	}

	l.InstallmentDetails = pq.Int64Array(installmentsID)

	for _, v := range cicilanbungas {
		l.TotalLoan += v + cicilanpokok
	}

	err = syncInstallment(installments)

	return err
}

func syncInstallment(is []Installment) error {
	type KafkaModelPayload struct {
		Payload interface{} `json:"payload"`
	}
	payload := KafkaModelPayload{
		Payload: is,
	}

	jMarshal, _ := json.Marshal(payload)

	config := sarama.NewConfig()
	config.ClientID = os.Getenv("KAFKA_CLIENT_ID")
	config.Net.SASL.Enable, _ = strconv.ParseBool(os.Getenv("KAFKA_SASL_ENABLE"))
	if config.Net.SASL.Enable {
		config.Net.SASL.User = os.Getenv("KAFKA_SASL_USER")
		config.Net.SASL.Password = os.Getenv("KAFKA_SASL_PASSWORD")
	}
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Consumer.Return.Errors = true
	topic := os.Getenv("KAFKA_PRODUCER_TOPIC")

	kafkaProducer, err := sarama.NewAsyncProducer([]string{os.Getenv("KAFKA_HOST")}, config)
	if err != nil {
		return err
	}
	defer kafkaProducer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder("installment_bulk:" + string(jMarshal)),
	}

	select {
	case kafkaProducer.Input() <- msg:
		log.Printf("Produced topic : %s", topic)
	case err := <-kafkaProducer.Errors():
		log.Printf("Fail producing topic : %s error : %v", topic, err)
	}

	return nil
}

// Create func
func (l *Loan) Create() error {
	return basemodel.Create(&l)
}

// FirstOrCreate func
func (l *Loan) FirstOrCreate() (err error) {
	return basemodel.FirstOrCreate(&l)
}

// Save func
func (l *Loan) Save() error {
	return basemodel.Save(&l)
}

// SaveNoKafka func
func (l *Loan) SaveNoKafka() error {
	err := basemodel.Save(&l)
	if err != nil {
		return err
	}

	return err
}

// Delete func
func (l *Loan) Delete() error {
	return basemodel.Delete(&l)
}

// FindbyID func
func (l *Loan) FindbyID(id uint64) error {
	return basemodel.FindbyID(&l, id)
}

// FilterSearchSingle func
func (l *Loan) FilterSearchSingle(filter interface{}) error {
	return basemodel.SingleFindFilter(&l, filter)
}

// PagedFilterSearch func
func (l *Loan) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	loans := []Loan{}
	var orders []string
	var sorts []string

	return basemodel.PagedFindFilter(&loans, page, rows, orders, sorts, filter)
}
