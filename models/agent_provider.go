package models

import (
	"asira_borrower/asira"

	"github.com/ayannahindonesia/basemodel"
)

// AgentProvider model
type AgentProvider struct {
	basemodel.BaseModel
	Name    string `json:"name" gorm:"column:name"`
	PIC     string `json:"pic" gorm:"column:pic"`
	Phone   string `json:"phone" gorm:"column:phone"`
	Address string `json:"address" gorm:"column:address"`
	Status  string `json:"status" gorm:"column:status"`
}

// Create new
func (model *AgentProvider) Create() error {
	return basemodel.Create(&model)
}

// BeforeSave gorm callback
func (model *AgentProvider) BeforeSave() error {
	if model.Status == "inactive" {
		deactivateAgents(model.ID)
	}
	return nil
}

// Save update
func (model *AgentProvider) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate saves, or create if not exist
func (model *AgentProvider) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete model
func (model *AgentProvider) Delete() error {
	return basemodel.Delete(&model)
}

// FindbyID func
func (model *AgentProvider) FindbyID(id uint64) error {
	return basemodel.FindbyID(&model, id)
}

// PagedFilterSearch paged list
func (model *AgentProvider) PagedFilterSearch(page int, rows int, order []string, sorts []string, filter interface{}) (basemodel.PagedFindResult, error) {
	agentProviders := []AgentProvider{}

	return basemodel.PagedFindFilter(&agentProviders, page, rows, order, sorts, filter)
}

func deactivateAgents(providerID uint64) {
	db := asira.App.DB

	db.Model(&Agent{}).Where("agent_provider = ?", providerID).Update("status", "inactive")
}
