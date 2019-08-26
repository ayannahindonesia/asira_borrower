package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
)

type (
	ServiceProduct struct {
		BaseModel
		Name            string         `json:"name" gorm:"column:name"`
		MinTimeSpan     int            `json:"min_timespan" gorm:"column:min_timespan"`
		MaxTimeSpan     int            `json:"max_timespan" gorm:"column:max_timespan"`
		Interest        int            `json:"interest" gorm:"column:interest"`
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

func (p *ServiceProduct) Create() (*ServiceProduct, error) {
	err := Create(&p)
	return p, err
}

func (p *ServiceProduct) Save() (*ServiceProduct, error) {
	err := Save(&p)
	return p, err
}

func (p *ServiceProduct) Delete() (*ServiceProduct, error) {
	err := Delete(&p)
	return p, err
}

func (p *ServiceProduct) FindbyID(id int) (*ServiceProduct, error) {
	err := FindbyID(&p, id)
	return p, err
}

func (p *ServiceProduct) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	product := []ServiceProduct{}
	result, err = PagedFilterSearch(&product, page, rows, orderby, sort, filter)

	return result, err
}

func (p *ServiceProduct) FilterSearch(filter interface{}) (SearchResult, error) {
	product := []ServiceProduct{}
	result, err := FilterSearch(&product, filter)
	return result, err
}
