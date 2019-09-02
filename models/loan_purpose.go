package models

import "time"

type LoanPurpose struct {
	BaseModel
	DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time"`
	Name        string    `json:"name" gorm:"column:name"`
	Status      string    `json:"status" gorm:"column:status"`
}

func (l *LoanPurpose) Create() (err error) {
	err = Create(&l)
	return err
}

func (l *LoanPurpose) Save() (err error) {
	err = Save(&l)
	return err
}

func (l *LoanPurpose) Delete() (err error) {
	l.DeletedTime = time.Now()
	err = Save(&l)

	return err
}
