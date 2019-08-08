package models

import (
	"database/sql"
	"time"

	guuid "github.com/google/uuid"
)

type (
	Uuid_Reset_Password struct {
		UUID        string        `json:"uuid" sql:"DEFAULT:NULL" gorm:"primary_key,column:uuid"`
		CreatedTime time.Time     `json:"created_time" gorm:"column:created_time" sql:"DEFAULT:current_timestamp"`
		UpdatedTime time.Time     `json:"updated_time" gorm:"column:updated_time" sql:"DEFAULT:current_timestamp"`
		Borrower    sql.NullInt64 `json:"borrower" gorm:"column:borrower" sql:"DEFAULT:NULL"`
		Expired     time.Time     `json:"expired" gorm:"column:expired"`
		Used        bool          `json:"used" gorm:"column:used;type:boolean" sql:"DEFAULT:FALSE"`
	}
)

// gorm callback hook
func (i *Uuid_Reset_Password) BeforeCreate() (err error) {
	id := guuid.New()
	i.UUID = id.String()

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
