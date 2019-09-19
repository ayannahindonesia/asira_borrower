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

func (i *Client_config) Create() (*Client_config, error) {
	err := Create(&i)
	return i, err
}

// gorm callback hook
func (i *Client_config) BeforeSave() (err error) {
	return nil
}

func (i *Client_config) Save() (*Client_config, error) {
	err := Save(&i)
	return i, err
}

func (l *Client_config) FilterSearchSingle(filter interface{}) (*Client_config, error) {
	err := FilterSearchSingle(&l, filter)
	return l, err
}

func (i *Client_config) Delete() (*Client_config, error) {
	err := Delete(&i)
	return i, err
}
