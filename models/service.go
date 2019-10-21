package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Service struct {
		basemodel.BaseModel
		DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time"`
		Name        string    `json:"name" gorm:"column:name;type:varchar(255)"`
		ImageID     uint64    `json:"image_id" gorm:"column:image_id"`
		Status      string    `json:"status" gorm:"column:status;type:varchar(255)"`
	}
)

func (model *Service) Create() error {
	err := basemodel.Create(&model)
	return err
}

func (model *Service) Save() error {
	err := basemodel.Save(&model)
	return err
}

func (model *Service) FirstOrCreate() (err error) {
	err = basemodel.FirstOrCreate(&model)
	return nil
}

func (model *Service) Delete() error {
	err := basemodel.Delete(&model)
	return err
}

func (model *Service) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

func (model *Service) FindFilter(order []string, sort []string, limit int, offset int, filter interface{}) (result interface{}, err error) {
	bank_type := []Service{}
	result, err = basemodel.FindFilter(&bank_type, order, sort, limit, offset, filter)
	return result, err
}

func (model *Service) PagedFindFilter(page int, rows int, orders []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	models := []Service{}
	result, err = basemodel.PagedFindFilter(&models, 0, 0, orders, sort, filter)
	return result, err
}

func (model *Service) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&model, filter)
	return err
}
