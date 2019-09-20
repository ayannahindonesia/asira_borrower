package models

import (
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	InternalRoles struct {
		basemodel.BaseModel
		Name        string `json:"name" gorm:"column:name"`
		Description string `json:"description" gorm:"column:description"`
		Status      bool   `json:"status" gorm:"column:status;type:boolean" sql:"DEFAULT:FALSE"`
		System      string `json:"system" gorm:"column:system"`
	}
)

func (b *InternalRoles) Create() error {
	err := basemodel.Create(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role")
	return err
}

func (b *InternalRoles) Save() error {
	err := basemodel.Save(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role")
	return err
}

func (b *InternalRoles) Delete() error {
	err := basemodel.Delete(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role_delete")
	return err
}

func (b *InternalRoles) FindbyID(id int) error {
	err := basemodel.FindbyID(&b, id)
	return err
}

func (b *InternalRoles) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&b, filter)
	return err
}

func (b *InternalRoles) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	internal := []InternalRoles{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&internal, page, rows, order, sorts, filter)

	return result, err
}

func (b *InternalRoles) FilterSearch(filter interface{}) (result basemodel.PagedFindResult, err error) {
	internal := []InternalRoles{}
	var orders []string
	var sort []string
	result, err = basemodel.PagedFindFilter(&internal, 0, 0, orders, sort, filter)
	return result, err
}
