package models

import (
	"database/sql"
	"time"

	"github.com/ayannahindonesia/basemodel"
)

type (
	Borrower struct {
		basemodel.BaseModel
		SuspendedTime        time.Time     `json:"suspended_time" gorm:"column:suspended_time"`
		Fullname             string        `json:"fullname" gorm:"column:fullname;type:varchar(255);not_null"`
		Nickname             string        `json:"nickname" gorm:"column:nickname"`
		Gender               string        `json:"gender" gorm:"column:gender;type:varchar(1);not null"`
		IdCardNumber         string        `json:"idcard_number" gorm:"column:idcard_number;type:varchar(255);unique;not null"`
		IdCardImage          string        `json:"idcard_image" gorm:"column:idcard_image;type:varchar(255)"`
		TaxIDnumber          string        `json:"taxid_number" gorm:"column:taxid_number;type:varchar(255)"`
		TaxIDImage           string        `json:"taxid_image" gorm:"column:taxid_image;type:varchar(255)"`
		Nationality          string        `json:"nationality" gorm:"column:nationality"`
		Email                string        `json:"email" gorm:"column:email;type:varchar(255);unique"`
		Birthday             time.Time     `json:"birthday" gorm:"column:birthday;not null"`
		Birthplace           string        `json:"birthplace" gorm:"column:birthplace;type:varchar(255);not null"`
		LastEducation        string        `json:"last_education" gorm:"column:last_education;type:varchar(255);not null"`
		MotherName           string        `json:"mother_name" gorm:"column:mother_name;type:varchar(255);not null"`
		Phone                string        `json:"phone" gorm:"column:phone;type:varchar(255);unique;not null"`
		MarriedStatus        string        `json:"marriage_status" gorm:"column:marriage_status;type:varchar(255);not null"`
		SpouseName           string        `json:"spouse_name" gorm:"column:spouse_name;type:varchar(255)"`
		SpouseBirthday       time.Time     `json:"spouse_birthday" gorm:"column:spouse_birthday"`
		SpouseLastEducation  string        `json:"spouse_lasteducation" gorm:"column:spouse_lasteducation;type:varchar(255)"`
		Dependants           int           `json:"dependants,omitempty" gorm:"column:dependants;type:int" sql:"DEFAULT:0"`
		Address              string        `json:"address" gorm:"column:address;type:varchar(255);not null"`
		Province             string        `json:"province" gorm:"column:province;type:varchar(255);not null"`
		City                 string        `json:"city" gorm:"column:city;type:varchar(255);not null"`
		NeighbourAssociation string        `json:"neighbour_association" gorm:"column:neighbour_association;type:varchar(255);not null"`
		Hamlets              string        `json:"hamlets" gorm:"column:hamlets;type:varchar(255);not null"`
		HomePhoneNumber      string        `json:"home_phonenumber" gorm:"column:home_phonenumber;type:varchar(255)"`
		Subdistrict          string        `json:"subdistrict" gorm:"column:subdistrict;type:varchar(255)";not null`
		UrbanVillage         string        `json:"urban_village" gorm:"column:urban_village;type:varchar(255)";not null`
		HomeOwnership        string        `json:"home_ownership" gorm:"column:home_ownership;type:varchar(255);not null`
		LivedFor             int           `json:"lived_for" gorm:"column:lived_for;type:int;not null"`
		Occupation           string        `json:"occupation" gorm:"column:occupation;type:varchar(255);not null"`
		EmployeeID           string        `json:"employee_id" gorm:"column:employee_id;type:varchar(255)"`
		EmployerName         string        `json:"employer_name" gorm:"column:employer_name;type:varchar(255);not null"`
		EmployerAddress      string        `json:"employer_address" gorm:"column:employer_address;type:varchar(255);not null"`
		Department           string        `json:"department" gorm:"column:department;type:varchar(255);not null"`
		BeenWorkingFor       int           `json:"been_workingfor" gorm:"column:been_workingfor;type:int;not null"`
		DirectSuperior       string        `json:"direct_superiorname" gorm:"column:direct_superiorname;type:varchar(255)"`
		EmployerNumber       string        `json:"employer_number" gorm:"column:employer_number;type:varchar(255);not null"`
		MonthlyIncome        int           `json:"monthly_income" gorm:"column:monthly_income;type:int;not null"`
		OtherIncome          int           `json:"other_income" gorm:"column:other_income;type:int"`
		OtherIncomeSource    string        `json:"other_incomesource" gorm:"column:other_incomesource;type:varchar(255)"`
		FieldOfWork          string        `json:"field_of_work" gorm:"column:field_of_work;type:varchar(255);not null"`
		RelatedPersonName    string        `json:"related_personname" gorm:"column:related_personname;type:varchar(255);not null"`
		RelatedRelation      string        `json:"related_relation" gorm:"column:related_relation;type:varchar(255);not null"`
		RelatedPhoneNumber   string        `json:"related_phonenumber" gorm:"column:related_phonenumber;type:varchar(255);not null"`
		RelatedHomePhone     string        `json:"related_homenumber" gorm:"column:related_homenumber;type:varchar(255)"`
		RelatedAddress       string        `json:"related_address" gorm:"column:related_address;type:text"`
		Bank                 sql.NullInt64 `json:"bank" gorm:"column:bank" sql:"DEFAULT:NULL"`
		BankAccountNumber    string        `json:"bank_accountnumber" gorm:"column:bank_accountnumber"`
		OTPverified          bool          `json:"otp_verified" gorm:"column:otp_verified;type:boolean" sql:"DEFAULT:FALSE"`
		AgentReferral        sql.NullInt64 `json:"agent_referral" gorm:"column:agent_referral"`
		Status               string        `json:"status" gorm:"column:status;default:'active'"`
		NthLoans             int           `json:"nth_loans" gorm:"-"`
	}
)

//Create just create row
func (b *Borrower) Create() error {
	err := basemodel.Create(&b)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(b, "borrower")
	return err
}

func (b *Borrower) Save() error {
	err := basemodel.Save(&b)
	if err != nil {
		return err
	}

	err = KafkaSubmitModel(b, "borrower")
	return err
}

func (b *Borrower) Suspend() error {
	b.SuspendedTime = time.Now()
	err := basemodel.Save(&b)

	return err
}

func (b *Borrower) Unsuspend() error {
	b.SuspendedTime = time.Time{}
	err := basemodel.Save(&b)

	return err
}

func (b *Borrower) FindbyID(id int) error {
	err := basemodel.FindbyID(&b, id)
	return err
}

func (model *Borrower) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&model, filter)
	return err
}

func (b *Borrower) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	borrowers := []Borrower{}
	var orders []string
	var sorts []string
	result, err = basemodel.PagedFindFilter(&borrowers, page, rows, orders, sorts, filter)

	return result, err
}
