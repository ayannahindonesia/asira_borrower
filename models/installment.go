package models

import (
	"time"

	"github.com/ayannahindonesia/basemodel"
)

// Installment details
type Installment struct {
	basemodel.BaseModel
	Period          int        `json:"period" gorm:"column:period"`
	LoanPayment     float64    `json:"loan_payment" gorm:"column:loan_payment"`
	InterestPayment float64    `json:"interest_payment" gorm:"column:interest_payment"`
	PaidDate        *time.Time `json:"paid_date" gorm:"column:paid_date"`
	PaidStatus      bool       `json:"paid_status" gorm:"column:paid_status" sql:"DEFAULT:FALSE"`
	Underpayment    float64    `json:"underpayment" gorm:"column:underpayment"`
	Note            string     `json:"note" gorm:"column:note"`
}

// Create func
func (model *Installment) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Installment) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate create if not exist, or skip if exist
func (model *Installment) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Installment) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Installment) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *Installment) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter func
func (model *Installment) PagedFindFilter(page int, rows int, orderby []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	lists := []Installment{}

	return basemodel.PagedFindFilter(&lists, page, rows, orderby, sort, filter)
}
