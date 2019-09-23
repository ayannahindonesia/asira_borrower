package models

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type (
	Loan struct {
		BaseModel
		DeletedTime      time.Time      `json:"deleted_time" gorm:"column:deleted_time"`
		Owner            sql.NullInt64  `json:"owner" gorm:"column:owner;foreignkey"`
		Status           string         `json:"status" gorm:"column:status;type:varchar(255)" sql:"DEFAULT:'processing'"`
		LoanAmount       float64        `json:"loan_amount" gorm:"column:loan_amount;type:int;not null"`
		Installment      int            `json:"installment" gorm:"column:installment;type:int;not null"` // plan of how long loan to be paid
		Fees             postgres.Jsonb `json:"fees" gorm:"column:fees;type:jsonb"`
		Interest         float64        `json:"interest" gorm:"column:interest;type:int;not null"`
		TotalLoan        float64        `json:"total_loan" gorm:"column:total_loan;type:int;not null"`
		DisburseAmount   float64        `json:"disburse_amount" gorm:"column:disburse_amount;type:int;not null"`
		DueDate          time.Time      `json:"due_date" gorm:"column:due_date"`
		LayawayPlan      float64        `json:"layaway_plan" gorm:"column:layaway_plan;type:int;not null"` // how much borrower will pay per month
		Product          uint64         `json:"product" gorm:"column:product;foreignkey"`
		Service          uint64         `json:"service" gorm:"column:service;foreignkey"`
		LoanIntention    string         `json:"loan_intention" gorm:"column:loan_intention;type:varchar(255);not null"`
		IntentionDetails string         `json:"intention_details" gorm:"column:intention_details;type:text;not null"`
		BorrowerInfo     postgres.Jsonb `json:"borrower_info" gorm:"column:borrower_info;type:jsonb"`
		OTPverified      bool           `json:"otp_verified" gorm:"column:otp_verified;type:boolean" sql:"DEFAULT:FALSE"`
	}

	LoanFee struct { // temporary hardcoded
		Description string `json:"description"`
		Amount      string `json:"amount"`
	}
	LoanFees []LoanFee
)

// gorm callback hook
func (l *Loan) BeforeCreate() (err error) {
	borrower := Borrower{}
	_, err = borrower.FindbyID(int(l.Owner.Int64))
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
		return err
	}

	return nil
}

func (l *Loan) SetProductLoanReferences() (err error) {
	product := ServiceProduct{}
	_, err = product.FindbyID(int(l.Product))
	if err != nil {
		return err
	}

	l.Fees = product.Fees
	l.Interest = product.Interest

	return nil
}

func (l *Loan) Calculate() (err error) {
	// calculate total loan
	var (
		totalfee float64
		fees     LoanFees
		owner    Borrower
		bank     Bank
		product  ServiceProduct
	)

	owner.FindbyID(int(l.Owner.Int64))
	bank.FindbyID(int(owner.Bank.Int64))
	product.FindbyID(int(l.Product))

	json.Unmarshal(l.Fees.RawMessage, &fees)

	var fee float64
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

		totalfee += fee
	}
	interest := (l.Interest / 100) * l.LoanAmount
	l.DisburseAmount = l.LoanAmount

	switch bank.AdminFeeSetup {
	case "potong_plafon":
		l.DisburseAmount = l.DisburseAmount - totalfee
		break
		// case "beban_plafon":
		// 	l.DisburseAmount = l.DisburseAmount + totalfee
		// 	break
	}

	var asnFee float64
	if strings.ContainsAny(product.ASN_Fee, "%") {
		asnFeeString := strings.TrimFunc(product.ASN_Fee, func(r rune) bool {
			return !unicode.IsNumber(r)
		})
		f, _ := strconv.Atoi(asnFeeString)
		asnFee = (float64(f) / 100) * l.LoanAmount
	} else {
		f, _ := strconv.Atoi(product.ASN_Fee)
		asnFee = float64(f)
	}

	switch bank.ConvinienceFeeSetup {
	case "potong_plafon":
		l.DisburseAmount = l.DisburseAmount - asnFee
		l.TotalLoan = l.LoanAmount + interest
		break
	case "beban_plafon":
		l.TotalLoan = l.LoanAmount + interest + asnFee + totalfee
		break
	}

	// calculate layaway plan
	l.LayawayPlan = l.TotalLoan / float64(l.Installment)

	return nil
}

func (l *Loan) Create() (*Loan, error) {
	err := Create(&l)
	return l, err
}

func (l *Loan) Save() (*Loan, error) {
	err := Save(&l)
	if err != nil {
		return nil, err
	}

	if l.OTPverified {
		err = KafkaSubmitModel(l, "loan")
	}
	return l, err
}

func (l *Loan) Delete() (*Loan, error) {
	l.DeletedTime = time.Now()
	err := Save(&l)

	return l, err
}

func (l *Loan) FindbyID(id int) (*Loan, error) {
	err := FindbyID(&l, id)
	return l, err
}

func (l *Loan) FilterSearchSingle(filter interface{}) (*Loan, error) {
	err := FilterSearchSingle(&l, filter)
	return l, err
}

func (l *Loan) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	loans := []Loan{}
	result, err = PagedFilterSearch(&loans, page, rows, orderby, sort, filter)

	return result, err
}
