package models

import (
	"time"

	"github.com/ayannahindonesia/basemodel"
)

type LoanPurpose struct {
	basemodel.BaseModel
	DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time"`
	Name        string    `json:"name" gorm:"column:name"`
	Status      string    `json:"status" gorm:"column:status"`
}

func (l *LoanPurpose) Create() (err error) {
	err = basemodel.Create(&l)
	return err
}

func (l *LoanPurpose) Save() (err error) {
	err = basemodel.Save(&l)
	return err
}

func (l *LoanPurpose) Delete() (err error) {
	l.DeletedTime = time.Now()
	err = basemodel.Save(&l)

	return err
}

func (l *LoanPurpose) FindbyID(id int) (err error) {
	err = basemodel.FindbyID(&l, id)
	return err
}

func (l *LoanPurpose) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&l, filter)
	return err
}

func (l *LoanPurpose) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	loan_purposes := []LoanPurpose{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&loan_purposes, page, rows, order, sorts, filter)

	return result, err
}
