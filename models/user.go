package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
	"golang.org/x/crypto/bcrypt"
)

//User model for table users
type User struct {
	basemodel.BaseModel
	DeletedTime time.Time `json:"deleted_time" gorm:"column:deleted_time"`
	Borrower    uint64    `json:"borrower" gorm:"column:borrower"`
	Password    string    `json:"password" gorm:"column:password"`
	FCMToken    string    `json:"fcm_token" gorm:"column:fcm_token;type:varchar(255)"`
}

// BeforeCreate gorm callback
func (model *User) BeforeCreate() (err error) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	model.Password = string(passwordByte)
	return nil
}

// Create new User
func (model *User) Create() error {
	err := basemodel.Create(&model)

	return err
}

// Save update User
func (model *User) Save() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}

	return err
}

// Delete User
func (model *User) Delete() error {
	err := basemodel.Delete(&model)
	if err != nil {
		return err
	}
	return err
}

// FindbyID find User with id
func (model *User) FindbyID(id int) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

// FilterSearchSingle search using filter and return last
func (model *User) FindbyBorrowerID(borrowerID uint64) error {
	type Filter struct {
		Borrower uint64 `json:"borrower_fk"`
	}
	err := basemodel.SingleFindFilter(&model, &Filter{
		Borrower: borrowerID,
	})
	return err
}

// FilterSearchSingle search using filter and return last
func (model *User) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&model, filter)
	return err
}

// PagedFilterSearch search using filter and return with pagination format
func (model *User) PagedFilterSearch(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	Users := []User{}
	result, err = basemodel.PagedFindFilter(&Users, page, rows, order, sort, filter)

	return result, err
}
