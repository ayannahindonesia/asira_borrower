package handlers

import (
	"asira_borrower/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func AgentRegisterBorrower(c echo.Context) error {
	defer c.Request().Body.Close()
	type (
		Register struct {
			Fullname             string    `json:"fullname"`
			Nickname             string    `json:"nickname"`
			Gender               string    `json:"gender" `
			IdCardNumber         string    `json:"idcard_number" `
			IdCardImage          string    `json:"idcard_image"`
			TaxIDImage           string    `json:"taxid_image"`
			TaxIDnumber          string    `json:"taxid_number"`
			Nationality          string    `json:"nationality"`
			Email                string    `json:"email"`
			Birthday             time.Time `json:"birthday"`
			Birthplace           string    `json:"birthplace"`
			LastEducation        string    `json:"last_education"`
			MotherName           string    `json:"mother_name"`
			Phone                string    `json:"phone"`
			MarriedStatus        string    `json:"marriage_status"`
			SpouseName           string    `json:"spouse_name"`
			SpouseBirthday       time.Time `json:"spouse_birthday"`
			SpouseLastEducation  string    `json:"spouse_lasteducation"`
			Dependants           int       `json:"dependants,omitempty"`
			Address              string    `json:"address"`
			Province             string    `json:"province"`
			City                 string    `json:"city"`
			NeighbourAssociation string    `json:"neighbour_association"`
			Hamlets              string    `json:"hamlets"`
			HomePhoneNumber      string    `json:"home_phonenumber"`
			Subdistrict          string    `json:"subdistrict"`
			UrbanVillage         string    `json:"urban_village"`
			HomeOwnership        string    `json:"home_ownership"`
			LivedFor             int       `json:"lived_for"`
			Occupation           string    `json:"occupation"`
			EmployeeID           string    `json:"employee_id"`
			EmployerName         string    `json:"employer_name"`
			EmployerAddress      string    `json:"employer_address"`
			Department           string    `json:"department"`
			BeenWorkingFor       int       `json:"been_workingfor"`
			DirectSuperior       string    `json:"direct_superiorname"`
			EmployerNumber       string    `json:"employer_number"`
			MonthlyIncome        int       `json:"monthly_income"`
			OtherIncome          int       `json:"other_income"`
			OtherIncomeSource    string    `json:"other_incomesource"`
			FieldOfWork          string    `json:"field_of_work"`
			RelatedPersonName    string    `json:"related_personname"`
			RelatedRelation      string    `json:"related_relation"`
			RelatedPhoneNumber   string    `json:"related_phonenumber"`
			RelatedHomePhone     string    `json:"related_homenumber"`
			RelatedAddress       string    `json:"related_address"`
			Bank                 int64     `json:"bank"`
			BankAccountNumber    string    `json:"bank_accountnumber"`
		}
	)
	register := Register{}

	payloadRules := govalidator.MapData{
		"fullname":              []string{"required"},
		"nickname":              []string{},
		"gender":                []string{"required"},
		"idcard_number":         []string{"required", "unique:agent_borrowers,idcard_number"},
		"taxid_number":          []string{"unique:agent_borrowers,taxid_number"},
		"nationality":           []string{},
		"email":                 []string{"email", "unique:agent_borrowers,email"},
		"birthday":              []string{"date"},
		"birthplace":            []string{"required"},
		"last_education":        []string{"required"},
		"mother_name":           []string{"required"},
		"phone":                 []string{"id_phonenumber", "unique:agent_borrowers,phone"},
		"marriage_status":       []string{"required"},
		"spouse_name":           []string{},
		"spouse_birthday":       []string{"date"},
		"spouse_lasteducation":  []string{},
		"dependants":            []string{},
		"address":               []string{"required"},
		"province":              []string{"required"},
		"city":                  []string{"required"},
		"neighbour_association": []string{"required"},
		"hamlets":               []string{"required"},
		"home_phonenumber":      []string{},
		"subdistrict":           []string{"required"},
		"urban_village":         []string{"required"},
		"home_ownership":        []string{"required"},
		"lived_for":             []string{"required"},
		"occupation":            []string{"required"},
		"employee_id":           []string{},
		"employer_name":         []string{"required"},
		"employer_address":      []string{"required"},
		"department":            []string{"required"},
		"been_workingfor":       []string{"required"},
		"direct_superiorname":   []string{},
		"employer_number":       []string{"required"},
		"monthly_income":        []string{"required"},
		"other_income":          []string{},
		"other_incomesource":    []string{},
		"field_of_work":         []string{"required"},
		"related_personname":    []string{"required"},
		"related_relation":      []string{"required"},
		"related_phonenumber":   []string{"required"},
		"related_homenumber":    []string{},
		"bank":                  []string{"required"},
		"bank_accountnumber":    []string{"unique:agent_borrowers,bank_accountnumber"},
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentModel := models.Agent{}
	agentID, _ := strconv.ParseInt(claims["jti"].(string), 10, 64)
	err := agentModel.FindbyID(int(agentID))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun Agen tidak ditemukan")
	}

	//validate
	validate := validateRequestPayload(c, payloadRules, &register)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek bank di dalam list bank agent atau tidak
	validBank := false
	for _, val := range agentModel.Banks {
		if register.Bank == val {
			validBank = true
			break
		}
	}
	if !validBank {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Bank tidak terdaftar untuk agent")
	}

	IdCardImage := models.Image{
		Image_string: register.IdCardImage,
	}
	err = IdCardImage.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	TaxIdImage := models.Image{
		Image_string: register.TaxIDImage,
	}
	err = TaxIdImage.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}
	borrower := models.AgentBorrower{
		AgentID: sql.NullInt64{
			Int64: agentID,
			Valid: true,
		},
		IdCardImage: sql.NullInt64{
			Int64: int64(IdCardImage.ID),
			Valid: true,
		},
		TaxIDImage: sql.NullInt64{
			Int64: int64(TaxIdImage.ID),
			Valid: true,
		},
		Bank: sql.NullInt64{
			Int64: int64(register.Bank),
			Valid: true,
		},
	}
	r, err := json.Marshal(register)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}
	json.Unmarshal(r, &borrower)

	err = borrower.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	return c.JSON(http.StatusCreated, borrower)
}
