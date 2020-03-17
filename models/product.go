package models

import (
	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

type (
	Product struct {
		basemodel.BaseModel
		Name            string         `json:"name" gorm:"column:name;type:varchar(255)"`
		ServiceID       uint64         `json:"service_id" gorm:"column:service_id`
		MinTimeSpan     int            `json:"min_timespan" gorm:"column:min_timespan"`
		MaxTimeSpan     int            `json:"max_timespan" gorm:"column:max_timespan"`
		Interest        float64        `json:"interest" gorm:"column:interest"`
		InterestType    string         `json:"interest_type" gorm:"column:interest_type"`
		MinLoan         int            `json:"min_loan" gorm:"column:min_loan"`
		MaxLoan         int            `json:"max_loan" gorm:"column:max_loan"`
		Fees            postgres.Jsonb `json:"fees" gorm:"column:fees"`
		Collaterals     pq.StringArray `json:"collaterals" gorm:"column:collaterals"`
		FinancingSector pq.StringArray `json:"financing_sector" gorm:"column:financing_sector"`
		Assurance       string         `json:"assurance" gorm:"column:assurance"`
		Status          string         `json:"status" gorm:"column:status;type:varchar(255)"`
		Form            postgres.Jsonb `json:"form" gorm:"column:form;type:text"`
	}
)

func (model *Product) Create() error {
	err := basemodel.Create(&model)
	return err
}

func (model *Product) Save() error {
	err := basemodel.Save(&model)
	return err
}

func (model *Product) FirstOrCreate() (err error) {
	return basemodel.FirstOrCreate(&model)
}

func (model *Product) Delete() error {
	err := basemodel.Delete(&model)
	return err
}

func (model *Product) FindbyID(id uint64) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *Product) FilterSearch(filter interface{}) (result basemodel.PagedFindResult, err error) {
	product := []Product{}
	var orders []string
	var sort []string
	result, err = basemodel.PagedFindFilter(&product, 0, 0, orders, sort, filter)
	return result, err
}

func (model *Product) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&model, filter)
	return err
}
