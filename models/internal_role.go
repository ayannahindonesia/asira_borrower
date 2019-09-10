package models

type (
	InternalRoles struct {
		BaseModel
		Name        string `json:"name" gorm:"column:name"`
		Description string `json:"description" gorm:"column:description"`
		Status      bool   `json:"status" gorm:"column:status;type:boolean" sql:"DEFAULT:FALSE"`
		System      string `json:"system" gorm:"column:system"`
	}
)

func (b *InternalRoles) Create() (*InternalRoles, error) {
	err := Create(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role")
	return b, err
}

func (b *InternalRoles) Save() (*InternalRoles, error) {
	err := Save(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role")
	return b, err
}

func (b *InternalRoles) Delete() (*InternalRoles, error) {
	err := Delete(&b)
	if err != nil {
		return nil, err
	}

	err = KafkaSubmitModel(b, "internal_role_delete")
	return b, err
}

func (b *InternalRoles) FindbyID(id int) (*InternalRoles, error) {
	err := FindbyID(&b, id)
	return b, err
}

func (b *InternalRoles) FilterSearchSingle(filter interface{}) (*InternalRoles, error) {
	err := FilterSearchSingle(&b, filter)
	return b, err
}

func (b *InternalRoles) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result PagedSearchResult, err error) {
	internal := []InternalRoles{}
	result, err = PagedFilterSearch(&internal, page, rows, orderby, sort, filter)

	return result, err
}

func (b *InternalRoles) FilterSearch(filter interface{}) (SearchResult, error) {
	internal := []InternalRoles{}
	result, err := FilterSearch(&internal, filter)
	return result, err
}
