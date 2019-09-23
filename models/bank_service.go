package models

import (
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	BankService struct {
		basemodel.BaseModel
		Name    string `json:"name" gorm:"column:name"`
		ImageID int    `json:"image_id" gorm:"column:image_id"`
		Status  string `json:"status" gorm:"column:status"`
	}
)

func (b *BankService) Create() error {
	err := basemodel.Create(&b)
	return err
}

func (b *BankService) Save() error {
	err := basemodel.Save(&b)
	return err
}

func (b *BankService) Delete() error {
	err := basemodel.Delete(&b)
	return err
}

func (b *BankService) FindbyID(id int) error {
	err := basemodel.FindbyID(&b, id)
	return err
}

func (b *BankService) FilterSearch(filter interface{}) (result basemodel.PagedFindResult, err error) {
	bank_type := []BankService{}
	var orders []string
	var sort []string
	result, err = basemodel.PagedFindFilter(&bank_type, 0, 0, orders, sort, filter)
	return result, err
}

func (b *BankService) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&b, filter)
	return err
}
