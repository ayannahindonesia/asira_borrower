package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	ServiceProduct struct {
		basemodel.BaseModel
		Name            string         `json:"name" gorm:"column:name"`
		MinTimeSpan     int            `json:"min_timespan" gorm:"column:min_timespan"`
		MaxTimeSpan     int            `json:"max_timespan" gorm:"column:max_timespan"`
		Interest        float64        `json:"interest" gorm:"column:interest"`
		MinLoan         int            `json:"min_loan" gorm:"column:min_loan"`
		MaxLoan         int            `json:"max_loan" gorm:"column:max_loan"`
		Fees            postgres.Jsonb `json:"fees" gorm:"column:fees"`
		ASN_Fee         string         `json:"asn_fee" gorm:"column:asn_fee"`
		Service         int            `json:"service" gorm:"column:service"`
		Collaterals     postgres.Jsonb `json:"collaterals" gorm:"column:collaterals"`
		FinancingSector postgres.Jsonb `json:"financing_sector" gorm:"column:financing_sector"`
		Assurance       string         `json:"assurance" gorm:"column:assurance"`
		Status          string         `json:"status" gorm:"column:status"`
	}
)

func (p *ServiceProduct) Create() error {
	err := basemodel.Create(&p)
	return err
}

func (p *ServiceProduct) Save() error {
	err := basemodel.Save(&p)
	return err
}

func (p *ServiceProduct) Delete() error {
	err := basemodel.Delete(&p)
	return err
}

func (p *ServiceProduct) FindbyID(id int) error {
	err := basemodel.FindbyID(&p, id)
	return err
}

func (p *ServiceProduct) FilterSearch(filter interface{}) (result basemodel.PagedFindResult, err error) {
	product := []ServiceProduct{}
	var orders []string
	var sort []string
	result, err = basemodel.PagedFindFilter(&product, 0, 0, orders, sort, filter)
	return result, err
}

func (p *ServiceProduct) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&p, filter)
	return err
}
