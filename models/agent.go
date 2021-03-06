package models

import (
	"database/sql"

	"github.com/ayannahindonesia/basemodel"
	"github.com/lib/pq"
)

// Agent main type
type Agent struct {
	basemodel.BaseModel
	Name          string        `json:"name" gorm:"column:name"`
	Username      string        `json:"username" gorm:"column:username"`
	Password      string        `json:"password" gorm:"column:password"`
	Email         string        `json:"email" gorm:"column:email"`
	Phone         string        `json:"phone" gorm:"column:phone"`
	Category      string        `json:"category" gorm:"column:category"`
	AgentProvider sql.NullInt64 `json:"agent_provider" gorm:"column:agent_provider"`
	Image         string        `json:"image" gorm:"column:image"`
	Banks         pq.Int64Array `json:"banks" gorm:"column:banks"`
	Status        string        `json:"status" gorm:"column:status"`
	FCMToken      string        `json:"fcm_token" gorm:"column:fcm_token;type:varchar(255)"`
}

// Create new agent
func (model *Agent) Create() error {
	err := basemodel.Create(&model)
	if err != nil {
		return err
	}
	return err
}

// Save update agent
func (model *Agent) Save() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}
	return err
}

func (model *Agent) FirstOrCreate() (err error) {
	return basemodel.FirstOrCreate(&model)
}

// Save update agent
func (model *Agent) SaveNoKafka() error {
	err := basemodel.Save(&model)
	if err != nil {
		return err
	}

	return err
}

// Delete agent
func (model *Agent) Delete() error {
	err := basemodel.Delete(&model)
	if err != nil {
		return err
	}
	return err
}

// FindbyID find agent with id
func (model *Agent) FindbyID(id uint64) error {
	err := basemodel.FindbyID(&model, id)
	return err
}

// FilterSearchSingle search using filter and return last
func (model *Agent) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&model, filter)
	return err
}

// PagedFilterSearch search using filter and return with pagination format
func (model *Agent) PagedFilterSearch(page int, rows int, order []string, sort []string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	agents := []Agent{}
	result, err = basemodel.PagedFindFilter(&agents, page, rows, order, sort, filter)

	return result, err
}

// checkBorrowerID search using filter and return last
func (model *Agent) CheckBorrowerOwnedByAgent(borrowerID uint64) bool {
	borrowerModel := Borrower{}
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return false
	}
	//cek agent id is correct
	if model.ID != uint64(borrowerModel.AgentReferral.Int64) {
		return false
	}
	return true
}
