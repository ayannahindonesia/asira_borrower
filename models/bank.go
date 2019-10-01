package models

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Bank struct {
		basemodel.BaseModel
		DeletedTime         time.Time      `json:"deleted_time" gorm:"column:deleted_time" sql:"DEFAULT:current_timestamp"`
		Name                string         `json:"name" gorm:"column:name;type:varchar(255)"`
		Type                int            `json:"type" gorm:"column:type;type:varchar(255)"`
		Address             string         `json:"address" gorm:"column:address;type:text"`
		Province            string         `json:"province" gorm:"column:province;type:varchar(255)"`
		City                string         `json:"city" gorm:"column:city;type:varchar(255)"`
		AdminFeeSetup       string         `json:"adminfee_setup" gorm:"column:adminfee_setup;type:varchar(255)"`
		ConvinienceFeeSetup string         `json:"convfee_setup" gorm:"column:convfee_setup;type:varchar(255)"`
		Services            postgres.Jsonb `json:"services" gorm:"column:services;type:jsonb"`
		Products            postgres.Jsonb `json:"products" gorm:"column:products;type:jsonb"`
		PIC                 string         `json:"pic" gorm:"column:pic;type:varchar(255)"`
		Phone               string         `json:"phone" gorm:"column:phone;type:varchar(255)"`
	}
)

func (model *Bank) Create() error {
	err := basemodel.Create(&model)
	return err
}

func (model *Bank) FirstOrCreate() (err error) {
	err = basemodel.FirstOrCreate(&model)
	return nil
}

func (model *Bank) Save() error {
	err := basemodel.Save(&model)
	return err
}

func (model *Bank) Delete() error {
	err := basemodel.Delete(&model)
	return err
}

func (model *Bank) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *Bank) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	banks := []Bank{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&banks, page, rows, order, sorts, filter)

	return result, err
}

func (model *Bank) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&model, filter)
	return err
}
