package models

import (
	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Image struct {
		basemodel.BaseModel
		Image_string string `json:"image_string" gorm:"column:image_string;type:text"`
	}
)

// gorm callback hook
func (i *Image) BeforeCreate() (err error) {
	return nil
}

func (i *Image) Create() error {
	err := basemodel.Create(&i)
	return err
}

// gorm callback hook
func (i *Image) BeforeSave() (err error) {
	return nil
}

func (i *Image) Save() error {
	err := basemodel.Save(&i)
	return err
}

func (i *Image) FindbyID(id int) error {
	err := basemodel.FindbyID(&i, id)
	return err
}
