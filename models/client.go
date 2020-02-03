package models

import (
	"github.com/ayannahindonesia/basemodel"
)

type (
	Client struct {
		basemodel.BaseModel
		Name   string `json:"name" gorm:"column:name"`
		Secret string `json:"secret" gorm:"column:secret"`
		Key    string `json:"key" gorm:"column:key"`
		Role   string `json:"role" gorm:"column:role"`
	}
)

// gorm callback hook
func (i *Client) BeforeCreate() (err error) {
	return nil
}

func (i *Client) Create() error {
	err := basemodel.Create(&i)
	return err
}

// gorm callback hook
func (i *Client) BeforeSave() (err error) {
	return nil
}

func (i *Client) Save() error {
	err := basemodel.Save(&i)
	return err
}

func (l *Client) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&l, filter)
	return err
}

func (i *Client) Delete() error {
	err := basemodel.Delete(&i)
	return err
}
