package models

import (
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Client_config struct {
		basemodel.BaseModel
		Name   string `json:"name" gorm:"column:name"`
		Secret string `json:"secret" gorm:"column:secret"`
		Key    string `json:"key" gorm:"column:key"`
		Role   string `json:"role" gorm:"column:role"`
	}
)

// gorm callback hook
func (i *Client_config) BeforeCreate() (err error) {
	return nil
}

func (i *Client_config) Create() error {
	err := basemodel.Create(&i)
	return err
}

// gorm callback hook
func (i *Client_config) BeforeSave() (err error) {
	return nil
}

func (i *Client_config) Save() error {
	err := basemodel.Save(&i)
	return err
}

func (l *Client_config) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&l, filter)
	return err
}

func (i *Client_config) Delete() error {
	err := basemodel.Delete(&i)
	return err
}
