package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"encoding/json"
	"fmt"
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

//RegisterBorrower register borrower personal
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
		"idcard_number":         []string{"required"},
		"taxid_number":          []string{},
		"nationality":           []string{},
		"email":                 []string{"required", "email"},
		"birthday":              []string{"date"},
		"birthplace":            []string{"required"},
		"last_education":        []string{"required"},
		"mother_name":           []string{"required"},
		"phone":                 []string{"required", "id_phonenumber"},
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
		"been_workingfor":       []string{},
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
		"bank_accountnumber":    []string{},
		"password":              []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &register)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	borrower := models.Borrower{}
	//marshalling data
	r, err := json.Marshal(register)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	//upload image id card
	IdCardImage, err := uploadImageS3Formatted("ktp", register.IdCardImage)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal : IDCardImage failed to upload")
	}

	//upload image tax card
	TaxIDImage, err := uploadImageS3Formatted("tax", register.TaxIDImage)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal : TaxIDImage failed to upload")
	}

	//create 1 user for 1 borrower (nasabah personal)
	user := models.User{}

	//search already exist Borrower registered by agent
	err = isBorrowerAlreadyRegistered(register.IdCardNumber)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Borrower personal sudah terdaftar sebelumnya")
	}

	//check manual fields if not unique
	var fields = map[string]string{
		"phone":              register.Phone,
		"email":              register.Email,
		"taxid_number":       register.TaxIDnumber,
		"bank_accountnumber": register.BankAccountNumber,
	}
	fieldsFound, err := checkUniqueFields(register.IdCardNumber, fields)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "data sudah ada sebelumnya : "+fieldsFound)
	}

	//create new personal borrower
	json.Unmarshal(r, &borrower)
	encryptPassphrase := asira.App.Config.GetString(fmt.Sprintf("%s.passphrase", asira.App.ENV))
	borrower.IdCardImage, err = encrypt(IdCardImage, encryptPassphrase)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi Id card gagal")
	}
	borrower.TaxIDImage, err = encrypt(TaxIDImage, encryptPassphrase)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi NPWP gagal")
	}

	borrower.Bank = sql.NullInt64{
		Int64: int64(register.Bank),
		Valid: true,
	}
	borrower.AgentReferral = sql.NullInt64{
		Int64: 0,
		Valid: true,
	}
	err = borrower.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	//save borrower_id to user entity and storing
	user.Borrower = borrower.ID
	user.Password = register.Password
	user.Create()

	return c.JSON(http.StatusCreated, borrower)
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

	message := fmt.Sprintf("Code OTP Registrasi anda adalah %s", otpCode)
	err := asira.App.Messaging.SendSMS(otpRequest.Phone, message)
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending otp")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Terkirim"})
}

func VerifyAccountOTP(c echo.Context) error {
	defer c.Request().Body.Close()

	otpVerify := VerifyAccountOTPverify{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	payloadRules := govalidator.MapData{
		"phone":    []string{"regex:^[0-9]+$", "required"},
		"otp_code": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &otpVerify)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	catenate := strconv.Itoa(int(borrowerID)) + otpVerify.Phone[len(otpVerify.Phone)-4:] // combine borrower id with last 4 digit of phone as counter
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

func updateAccountOTPstatus(borrowerID uint64) {
	modelBorrower := models.Borrower{}
	_ = modelBorrower.FindbyID(borrowerID)
	modelBorrower.OTPverified = true
	modelBorrower.Save()
}
