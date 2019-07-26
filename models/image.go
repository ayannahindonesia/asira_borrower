package models

type (
	Image struct {
		BaseModel
		Image_name string `json:"image_name" gorm:"column:image_name"`
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
