package models

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type (
	Bank struct {
		BaseModel
		DeletedTime time.Time      `json:"deleted_time" gorm:"column:deleted_time" sql:"DEFAULT:current_timestamp"`
		Name        string         `json:"name" gorm:"column:name;type:varchar(255)"`
		Type        int            `json:"type" gorm:"column:type;type:varchar(255)"`
		Address     string         `json:"address" gorm:"column:address;type:text"`
		Province    string         `json:"province" gorm:"column:province;type:varchar(255)"`
		City        string         `json:"city" gorm:"column:city;type:varchar(255)"`
		Services    postgres.Jsonb `json:"services" gorm:"column:services;type:jsonb"`
		Products    postgres.Jsonb `json:"products" gorm:"column:products;type:jsonb"`
		PIC         string         `json:"pic" gorm:"column:pic;type:varchar(255)"`
		Phone       string         `json:"phone" gorm:"column:phone;type:varchar(255)"`
	}
)

// gorm callback hook
func (b *Bank) BeforeCreate() (err error) {
	return nil
}

func (b *Bank) Create() (*Bank, error) {
	err := Create(&b)

	KafkaSubmitModel(b, "bank")

	return b, err
}

// gorm callback hook
func (b *Bank) BeforeSave() (err error) {
	return nil
}

func (b *Bank) Save() (*Bank, error) {
	err := Save(&b)
	return b, err
}

func (b *Bank) Delete() (*Bank, error) {
	err := Delete(&b)
	return b, err
}

func (b *Bank) FindbyID(id int) (*Bank, error) {
	err := FindbyID(&b, id)
	return b, err
}

func (b *Bank) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	bank_type := []Bank{}
	result, err = PagedFilterSearch(&bank_type, page, rows, orderby, sort, filter)

	return result, err
}
