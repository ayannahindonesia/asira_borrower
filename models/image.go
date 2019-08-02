package models

type (
	Image struct {
		BaseModel
		Image_string string `json:"image_string" gorm:"column:image_string;type:text"`
	}
)

// gorm callback hook
func (i *Image) BeforeCreate() (err error) {
	return nil
}

func (i *Image) Create() (*Image, error) {
	err := Create(&i)
	return i, err
}

// gorm callback hook
func (i *Image) BeforeSave() (err error) {
	return nil
}

func (i *Image) Save() (*Image, error) {
	err := Save(&i)
	return i, err
}

func (i *Image) FindbyID(id int) (*Image, error) {
	err := FindbyID(&i, id)
	return i, err
}
