package models

import (
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	BankService struct {
		basemodel.BaseModel
		Name    string `json:"name" gorm:"column:name"`
		BankID  uint64 `json:"bank_id" gorm:"column:bank_id"`
		ImageID uint64 `json:"image_id" gorm:"column:image_id"`
		Status  string `json:"status" gorm:"column:status"`
	}
)

func (model *BankService) Create() error {
	err := basemodel.Create(&model)
	return err
}

func (model *BankService) Save() error {
	err := basemodel.Save(&model)
	return err
}

func (model *BankService) FirstOrCreate() (err error) {
	err = basemodel.FirstOrCreate(&model)
	return nil
}

func (model *BankService) Delete() error {
	err := basemodel.Delete(&model)
	return err
}

func (model *BankService) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *BankService) FilterSearch(filter interface{}) (result basemodel.PagedFindResult, err error) {
	bank_type := []BankService{}
	var orders []string
	var sort []string
	result, err = basemodel.PagedFindFilter(&bank_type, 0, 0, orders, sort, filter)
	return result, err
}

func (model *BankService) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&model, filter)
	return err
}
