package models

import (
	"time"

	"github.com/ayannahindonesia/basemodel"
)

type (
	BankType struct {
		basemodel.BaseModel
		Name        string    `json:"name" gorm:"name"`
		Description string    `json:"description" gorm:"description"`
	}
)

func (b *BankType) Create() error {
	err := basemodel.Create(&b)
	return err
}

func (b *BankType) Save() error {
	err := basemodel.Save(&b)
	return err
}

func (model *BankType) FirstOrCreate() (err error) {
	err = basemodel.FirstOrCreate(&model)
	return nil
}

func (b *BankType) Delete() error {
	err := basemodel.Delete(&b)
	return err
}

func (b *BankType) FindbyID(id uint64) error {
	err := basemodel.FindbyID(&b, id)
	return err
}

func (b *BankType) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	bank_type := []BankType{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&bank_type, page, rows, order, sorts, filter)

	return result, err
}

func (b *BankType) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&b, filter)
	return err
}
