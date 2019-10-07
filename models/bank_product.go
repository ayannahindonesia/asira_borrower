package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	BankProduct struct {
		basemodel.BaseModel
		BankServiceID   uint64         `json:"bank_service_id"`
		Name            string         `json:"name" gorm:"column:name"`
		MinTimeSpan     int            `json:"min_timespan" gorm:"column:min_timespan"`
		MaxTimeSpan     int            `json:"max_timespan" gorm:"column:max_timespan"`
		Interest        float64        `json:"interest" gorm:"column:interest"`
		MinLoan         int            `json:"min_loan" gorm:"column:min_loan"`
		MaxLoan         int            `json:"max_loan" gorm:"column:max_loan"`
		Fees            postgres.Jsonb `json:"fees" gorm:"column:fees"`
		Collaterals     pq.StringArray `json:"collaterals" gorm:"column:collaterals"`
		FinancingSector pq.StringArray `json:"financing_sector" gorm:"column:financing_sector"`
		Assurance       string         `json:"assurance" gorm:"column:assurance"`
		Status          string         `json:"status" gorm:"column:status"`
	}
)

func (model *BankProduct) Create() error {
	err := basemodel.Create(&model)
	return err
}

func (model *BankProduct) Save() error {
	err := basemodel.Save(&model)
	return err
}

func (model *BankProduct) FirstOrCreate() (err error) {
	err = basemodel.FirstOrCreate(&model)
	return nil
}

func (model *BankProduct) Delete() error {
	err := basemodel.Delete(&model)
	return err
}

func (model *BankProduct) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *BankProduct) FilterSearch(filter interface{}) (result basemodel.PagedFindResult, err error) {
	product := []BankProduct{}
	var orders []string
	var sort []string
	result, err = basemodel.PagedFindFilter(&product, 0, 0, orders, sort, filter)
	return result, err
}

func (model *BankProduct) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&model, filter)
	return err
}
