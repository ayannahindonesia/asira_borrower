package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	VerifyAccountOTPrequest struct {
		Phone string `json:"phone"`
	}
	VerifyAccountOTPverify struct {
		VerifyAccountOTPrequest
		OTPcode string `json:"otp_code"`
	}
)

func RegisterBorrower(c echo.Context) error {
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
			Bank                 int       `json:"bank"`
			BankAccountNumber    string    `json:"bank_accountnumber"`
			Password             string    `json:"password"`
		}
	)
	register := Register{}
	payloadRules := govalidator.MapData{
		"fullname":              []string{"required"},
		"nickname":              []string{},
		"gender":                []string{"required"},
		"idcard_number":         []string{"required", "unique:borrowers,idcard_number"},
		"taxid_number":          []string{"unique:borrowers,taxid_number"},
		"nationality":           []string{},
		"email":                 []string{"email", "unique:borrowers,email"},
		"birthday":              []string{"date"},
		"birthplace":            []string{"required"},
		"last_education":        []string{"required"},
		"mother_name":           []string{"required"},
		"phone":                 []string{"required", "unique:borrowers,phone"},
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
		"bank":                  []string{},
		"bank_accountnumber":    []string{"unique:borrowers,bank_accountnumber"},
		"password":              []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &register)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}
	image := models.Image{
		Image_string: register.IdCardImage,
	}
	IdCardImage, err := image.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	image = models.Image{
		Image_string: register.TaxIDImage,
	}
	TaxIdImage, err := image.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}
	borrower := models.Borrower{
		IdCardImage: sql.NullInt64{
			Int64: int64(IdCardImage.BaseModel.ID),
			Valid: true,
		},
		TaxIDImage: sql.NullInt64{
			Int64: int64(TaxIdImage.BaseModel.ID),
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

	newBorrower, err := borrower.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	return c.JSON(http.StatusCreated, newBorrower)
}
func RequestOTPverifyAccount(c echo.Context) error {
	defer c.Request().Body.Close()

	otpRequest := VerifyAccountOTPrequest{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	payloadRules := govalidator.MapData{
		"phone": []string{"regex:^[0-9]+$", "required"},
	}

	validate := validateRequestPayload(c, payloadRules, &otpRequest)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	catenate := strconv.Itoa(borrowerID) + otpRequest.Phone[len(otpRequest.Phone)-4:] // combine borrower id with last 4 digit of phone as counter
	counter, _ := strconv.Atoi(catenate)
	otpCode := asira.App.OTP.HOTP.At(int(counter))
	log.Printf("generate otp code for phonenumber %s : %s", otpRequest.Phone, otpCode)

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Terkirim"})
}

func VerifyAccountOTP(c echo.Context) error {
	defer c.Request().Body.Close()

	otpVerify := VerifyAccountOTPverify{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	payloadRules := govalidator.MapData{
		"phone":    []string{"regex:^[0-9]+$", "required"},
		"otp_code": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &otpVerify)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	catenate := strconv.Itoa(borrowerID) + otpVerify.Phone[len(otpVerify.Phone)-4:] // combine borrower id with last 4 digit of phone as counter
	counter, _ := strconv.Atoi(catenate)
	if asira.App.OTP.HOTP.Verify(otpVerify.OTPcode, counter) {
		updateAccountOTPstatus(borrowerID)
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Verified"})
	}

	// bypass otp
	if asira.App.ENV == "development" && otpVerify.OTPcode == "888999" {
		updateAccountOTPstatus(borrowerID)
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Verified"})
	}

	return returnInvalidResponse(http.StatusBadRequest, "", "OTP salah")
}

func updateAccountOTPstatus(borrowerID int) {
	modelBorrower := models.Borrower{}
	borrower, _ := modelBorrower.FindbyID(borrowerID)
	borrower.OTPverified = true
	borrower.Save()
}
