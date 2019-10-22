package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Loan struct {
		basemodel.BaseModel
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
		LoanIntention    string         `json:"loan_intention" gorm:"column:loan_intention;type:varchar(255);not null"`
		IntentionDetails string         `json:"intention_details" gorm:"column:intention_details;type:text;not null"`
		BorrowerInfo     postgres.Jsonb `json:"borrower_info" gorm:"column:borrower_info;type:jsonb"`
		OTPverified      bool           `json:"otp_verified" gorm:"column:otp_verified;type:boolean" sql:"DEFAULT:FALSE"`
		DisburseDate     time.Time      `json:"disburse_date" gorm:"column:disburse_date"`
		DisburseStatus   string         `json:"disburse_status" gorm:"column:disburse_status" sql:"DEFAULT:'processing'"`
		RejectReason     string         `json:"reject_reason" gorm:"column:reject_reason"`
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
	err = borrower.FindbyID(int(l.Owner.Int64))
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
	product := Product{}
	err = product.FindbyID(int(l.Product))
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
		totalfee       float64
		fee            float64
		convenienceFee float64
		fees           LoanFees
		owner          Borrower
		bank           Bank
		product        Product
		parsedFees     LoanFees
	)

	owner.FindbyID(int(l.Owner.Int64))
	bank.FindbyID(int(owner.Bank.Int64))
	product.FindbyID(int(l.Product))

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
			Amount:      fmt.Sprint(fee),
		})

		if strings.ToLower(v.Description) == "convenience fee" {
			convenienceFee += fee
		} else {
			totalfee += fee
		}
	}
	// parse fees
	jMarshal, _ := json.Marshal(parsedFees)
	l.Fees = postgres.Jsonb{jMarshal}

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

	switch bank.ConvenienceFeeSetup {
	case "potong_plafon":
		l.DisburseAmount = l.DisburseAmount - convenienceFee
		l.TotalLoan = l.LoanAmount + interest
		break
	case "beban_plafon":
		l.TotalLoan = l.LoanAmount + interest + convenienceFee + totalfee
		break
	}

	// calculate layaway plan
	l.LayawayPlan = l.TotalLoan / float64(l.Installment)

	return nil
}

func (l *Loan) Create() error {
	err := basemodel.Create(&l)
	if err != nil {
		return err
	}

	if l.OTPverified {
		err = KafkaSubmitModel(l, "loan")
	}
	return err
}

func (l *Loan) Save() error {
	err := basemodel.Save(&l)
	if err != nil {
		return err
	}

	if l.OTPverified {
		err = KafkaSubmitModel(l, "loan")
	}
	return err
}

func (l *Loan) Delete() error {
	l.DeletedTime = time.Now()
	err := basemodel.Save(&l)

	return err
}

func (l *Loan) FindbyID(id int) error {
	err := basemodel.FindbyID(&l, id)
	return err
}

func (l *Loan) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&l, filter)
	return err
}

func (l *Loan) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	loans := []Loan{}
	var orders []string
	var sorts []string
	result, err = basemodel.PagedFindFilter(&loans, page, rows, orders, sorts, filter)

	return result, err
}
