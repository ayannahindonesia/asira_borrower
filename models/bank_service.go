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

func (b *BankService) Create() (*BankService, error) {
	err := Create(&b)
	return b, err
}

func (b *BankService) Save() (*BankService, error) {
	err := Save(&b)
	return b, err
}

func (b *BankService) Delete() (*BankService, error) {
	err := Delete(&b)
	return b, err
}

func (b *BankService) FindbyID(id int) (*BankService, error) {
	err := FindbyID(&b, id)
	return b, err
}

func (b *BankService) FilterSearchSingle(filter interface{}) (*BankService, error) {
	err := FilterSearchSingle(&b, filter)
	return b, err
}

func (b *BankService) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	bank_type := []BankService{}
	result, err = PagedFilterSearch(&bank_type, page, rows, orderby, sort, filter)

	return result, err
}

func (b *BankService) FilterSearch(filter interface{}) (SearchResult, error) {
	bank_type := []BankService{}
	result, err := FilterSearch(&bank_type, filter)
	return result, err
}
