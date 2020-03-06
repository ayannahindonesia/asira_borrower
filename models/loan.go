package models

import (
	"asira_borrower/custommodule/irate"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type (
	Loan struct {
		basemodel.BaseModel
		Borrower            uint64         `json:"borrower" gorm:"column:borrower;foreignkey"`
		Status              string         `json:"status" gorm:"column:status;type:varchar(255)" sql:"DEFAULT:'processing'"`
		LoanAmount          float64        `json:"loan_amount" gorm:"column:loan_amount;type:int;not null"`
		Installment         int            `json:"installment" gorm:"column:installment;type:int;not null"` // plan of how long loan to be paid
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
	}

	LoanFee struct {
		Description string `json:"description"`
		Amount      string `json:"amount"`
		FeeMethod   string `json:"fee_method"`
	}
	LoanFees []LoanFee
)

// gorm callback hook
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
		return err
	}

	return nil
}

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

func (l *Loan) Calculate() (err error) {
	// calculate total loan
	var (
		fee        float64
		fees       LoanFees
		borrower   Borrower
		bank       Bank
		product    Product
		parsedFees LoanFees
	)

	borrower.FindbyID(l.Borrower)
	bank.FindbyID(uint64(borrower.Bank.Int64))
	product.FindbyID(l.Product)

	l.CalculateInterest(product)
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
func (l *Loan) CalculateInterest(p Product) {
	switch p.InterestType {
	default:
		break
	case "flat":
		l.LayawayPlan, l.TotalLoan = irate.FLATANNUAL(l.Interest/100, l.LoanAmount, float64(l.Installment))
		break
	case "onetimepay":
		l.LayawayPlan, l.TotalLoan = irate.ONETIMEPAYMENT(l.Interest/100, l.LoanAmount, float64(l.Installment))
		break
	case "fixed":
		l.FixedInterestFormula()
		break
	case "efektif_menurun":
		l.EfektifMenurunFormula()
		break
	}
}

// FixedInterestFormula func
func (l *Loan) FixedInterestFormula() {
	rate := ((l.Interest / 100) / 12)
	pokok, bunga := irate.PIPMT(rate, 1, float64(l.Installment), -l.LoanAmount, 1)

	log.Println("pkok : %v \n bunga : %v", pokok, bunga)

	l.LayawayPlan = pokok + bunga
	l.TotalLoan = l.LayawayPlan * float64(l.Installment)
}

// EfektifMenurunFormula func
func (l *Loan) EfektifMenurunFormula() {
	plafon := l.LoanAmount
	cicilanpokok := l.LoanAmount / float64(l.Installment)
	var cicilanbungas []float64
	for i := 1; i <= l.Installment; i++ {
		bunga := plafon * (l.Interest / 100) / 12
		cicilanbungas = append(cicilanbungas, bunga)
		plafon -= cicilanpokok
	}
	log.Println("cek : %v \n cicilan pokok : %v", cicilanbungas, cicilanpokok)
	for _, v := range cicilanbungas {
		l.TotalLoan += v + cicilanpokok
	}
}

func (l *Loan) Create() error {
	err := basemodel.Create(&l)
	if err != nil {
		return err
	}

	// if l.OTPverified {
	// 	err = KafkaSubmitModel(l, "loan")
	// }
	return err
}

func (l *Loan) FirstOrCreate() (err error) {
	return basemodel.FirstOrCreate(&l)
}

func (l *Loan) Save() error {
	err := basemodel.Save(&l)
	if err != nil {
		return err
	}

	return err
}

func (l *Loan) SaveNoKafka() error {
	err := basemodel.Save(&l)
	if err != nil {
		return err
	}

	return err
}

func (l *Loan) Delete() error {
	err := basemodel.Delete(&l)
	return err
}

func (l *Loan) FindbyID(id uint64) error {
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
