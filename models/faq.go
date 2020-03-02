package models

import (
	"github.com/ayannahindonesia/basemodel"
)

// FAQ main type
type FAQ struct {
	basemodel.BaseModel
	Title       string `json:"title" gorm:"column:title"`
	Description string `json:"description" gorm:"column:description"`
}

// Create func
func (model *FAQ) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *FAQ) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate create if not exist, or skip if exist
func (model *FAQ) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete func
func (model *FAQ) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *FAQ) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// SingleFindFilter func
func (model *FAQ) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter func
func (model *FAQ) PagedFindFilter(page int, rows int, orderby []string, sort []string, filter interface{}) (basemodel.PagedFindResult, error) {
	loanpurposes := []FAQ{}

	return basemodel.PagedFindFilter(&loanpurposes, page, rows, orderby, sort, filter)
}
