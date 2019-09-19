package models

import (
	"database/sql"
	"time"

	basemodel "gitlab.com/asira-ayannah/basemodel"
)

type (
	Uuid_Reset_Password struct {
		basemodel.BaseModel
		UUID     string        `json:"uuid" sql:"DEFAULT:NULL" gorm:"primary_key,column:uuid"`
		Borrower sql.NullInt64 `json:"borrower" gorm:"column:borrower" sql:"DEFAULT:NULL"`
		Expired  time.Time     `json:"expired" gorm:"column:expired"`
		Used     bool          `json:"used" gorm:"column:used;type:boolean" sql:"DEFAULT:FALSE"`
	}
)

// gorm callback hook
func (i *Uuid_Reset_Password) BeforeCreate() (err error) {

	myDate := time.Now()
	i.Expired = myDate.AddDate(0, 0, 1)

	return nil
}

func (i *Uuid_Reset_Password) Create() (*Uuid_Reset_Password, error) {
	err := Create(&i)
	return i, err
}

// gorm callback hook
func (i *Uuid_Reset_Password) BeforeSave() (err error) {
	// i.Used = true
	return nil
}

func (i *Uuid_Reset_Password) Save() (*Uuid_Reset_Password, error) {
	err := Save(&i)
	return i, err
}

func (l *Uuid_Reset_Password) FilterSearchSingle(filter interface{}) (*Uuid_Reset_Password, error) {
	err := FilterSearchSingle(&l, filter)
	return l, err
}

func (i *Uuid_Reset_Password) Delete() (*Uuid_Reset_Password, error) {
	err := Delete(&i)
	return i, err
}
