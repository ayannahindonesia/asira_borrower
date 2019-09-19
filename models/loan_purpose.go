package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type LoanPurpose struct {
	basemodel.BaseModel
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

func (l *LoanPurpose) FindbyID(id int) (err error) {
	err = FindbyID(&l, id)
	return err
}

func (l *LoanPurpose) FilterSearchSingle(filter interface{}) (err error) {
	err = FilterSearchSingle(&l, filter)
	return err
}

func (l *LoanPurpose) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	loan_purposes := []LoanPurpose{}
	result, err = PagedFilterSearch(&loan_purposes, page, rows, orderby, sort, filter)

	return result, err
}
