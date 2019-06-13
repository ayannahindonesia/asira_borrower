package models

import (
	"kayacredit/kc"
	"time"

	"github.com/jinzhu/gorm"
)

type (
	BaseModel struct {
		ID          uint64    `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key,column:id"`
		CreatedTime time.Time `json:"created_time" gorm:"column:created_at" sql:"DEFAULT:current_timestamp"`
		UpdatedTime time.Time `json:"updated_time" gorm:"column:created_at" sql:"DEFAULT:current_timestamp"`
	}

	DBFunc func(tx *gorm.DB) error
)

// helper for inserting data using gorm.DB functions
func WithinTransaction(fn DBFunc) (err error) {
	tx := kc.App.DB.Begin()
	defer tx.Commit()
	err = fn(tx)

	return err
}

// inserts a row into db.
func Create(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		if !kc.App.DB.NewRecord(i) {
			return err
		}
		if err = tx.Create(i).Error; err != nil {
			tx.Rollback()
			return err
		}
		return err
	})
}

// Update row in db.
func Save(i interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		// check new object
		if kc.App.DB.NewRecord(i) {
			return err
		}
		if err = tx.Save(i).Error; err != nil {
			tx.Rollback() // rollback
			return err
		}
		return err
	})
}
