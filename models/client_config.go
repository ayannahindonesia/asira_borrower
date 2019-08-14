package models

import (
	guuid "github.com/google/uuid"
)

type (
	Client_config struct {
		BaseModel
		Name   string `json:"client_name" gorm:"primary_key,column:client_name"`
		Secret string `json:"secret" gorm:"column:secret"`
		Role   string `json:"role" gorm:"column:role"`
	}
)

// gorm callback hook
func (i *Client_config) BeforeCreate() (err error) {
	i.Secret = guuid.New().String()
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
