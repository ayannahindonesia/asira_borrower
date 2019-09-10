package models

type (
	Internal_Roles struct {
		BaseModel
		Name        string `json:"name" gorm:"column:name"`
		Description string `json:"description" gorm:"column:description"`
		Status      bool   `json:"status" gorm:"column:status;type:boolean" sql:"DEFAULT:FALSE"`
		System      string `json:"system" gorm:"column:system"`
	}
)

func (b *Internal_Roles) Create() (*Internal_Roles, error) {
	err := Create(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role")
	return b, err
}

func (b *Internal_Roles) Save() (*Internal_Roles, error) {
	err := Save(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role")
	return b, err
}

func (b *Internal_Roles) Delete() (*Internal_Roles, error) {
	err := Delete(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role_delete")
	return b, err
}

func (b *Internal_Roles) FindbyID(id int) (*Internal_Roles, error) {
	err := FindbyID(&b, id)
	return b, err
}

func (b *Internal_Roles) FilterSearchSingle(filter interface{}) (*Internal_Roles, error) {
	err := FilterSearchSingle(&b, filter)
	return b, err
}

func (b *Internal_Roles) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	internal := []Internal_Roles{}
	result, err = PagedFilterSearch(&internal, page, rows, orderby, sort, filter)

	return result, err
}

func (b *Internal_Roles) FilterSearch(filter interface{}) (SearchResult, error) {
	internal := []Internal_Roles{}
	result, err := FilterSearch(&internal, filter)
	return result, err
}
