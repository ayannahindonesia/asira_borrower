package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type (
	//Notification datatype
	Notification struct {
		basemodel.BaseModel
		ClientID    uint64 `json:"client_id" gorm:"column:client_id"`
		RecipientID uint64 `json:"recipient_id" gorm:"recipient_id"`
		Title       string `json:"title" gorm:"column:title"`
		MessageBody string `json:"message_body" gorm:"column:message_body"`
		//TODO: to get from client device
		FirebaseToken string    `json:"firebase_token" gorm:"column:firebase_token""`
		Topic         string    `json:"topic" gorm:"column:topic"`
		Response      string    `json:"response" gorm:"column:response"`
		SendTime      time.Time `json:"send_time" gorm:"column:send_time" sql:"DEFAULT:current_timestamp"`
	}
)

//BeforeCreate gorm callback hook
func (u *Notification) BeforeCreate() (err error) {
	return nil
}

//Create new Notification data
func (u *Notification) Create() error {
	err := basemodel.Create(&u)
	return err
}

//BeforeSave gorm callback hook
func (u *Notification) BeforeSave() (err error) {
	return nil
}

//Save / update data notification
func (u *Notification) Save() error {
	err := basemodel.Save(&u)
	return err
}

//FindbyID to search 1 row by ID
func (u *Notification) FindbyID(id int) error {
	err := basemodel.FindbyID(&u, id)
	return err
}

//FilterSearchSingle to search 1 row by filter (multiple fields)
func (u *Notification) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&u, filter)
	return err
}

//PagedFilterSearch FilterSearchSingle to search 1 row by filter (multiple fields) and paged properties
func (u *Notification) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	notif := []Notification{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&notif, page, rows, order, sorts, filter)

	return result, err
}
