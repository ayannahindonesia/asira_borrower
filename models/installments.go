package models

import (
	"github.com/ayannahindonesia/basemodel"
)

type Installments struct {
	basemodel.BaseModel
	LoanID          uint64  `json:"loan_id" gorm:"column:loan_id"`
	Period          int     `json:"period" gorm:"column:period"`
	LoanPayment     float64 `json:"loan_payment" gorm:"column:loan_payment"`
	InterestPayment float64 `json:"interest_payment" gorm:"column:interest_payment"`
}

// Create func
func (model *Installments) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Installments) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate create if not exist, or skip if exist
func (model *Installments) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *Installments) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *Installments) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *Installments) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter func
func (model *Installments) PagedFindFilter(page int, rows int, orderby []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	lists := []Installments{}

	return basemodel.PagedFindFilter(&lists, page, rows, orderby, sort, filter)
}
