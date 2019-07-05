package handlers

import (
	"asira/asira"
	"asira/models"
	"log"
	"net/http"
	"strconv"

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

	borrower := models.Borrower{}

	payloadRules := govalidator.MapData{
		"fullname":              []string{"required"},
		"gender":                []string{"required"},
		"idcard_number":         []string{"required", "unique:borrowers,idcard_number"},
		"taxid_number":          []string{"unique:borrowers,taxid_number"},
		"email":                 []string{"email", "unique:borrowers,email"},
		"birthday":              []string{"date"},
		"birthplace":            []string{"required"},
		"last_education":        []string{"required"},
		"mother_name":           []string{"required"},
		"phone":                 []string{"required"},
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

	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	newBorrower, err := borrower.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "create new borrower failed")
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

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Sent"})
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

	return returnInvalidResponse(http.StatusBadRequest, "", "Wrong OTP code")
}

func updateAccountOTPstatus(borrowerID int) {
	modelBorrower := models.Borrower{}
	borrower, _ := modelBorrower.FindbyID(borrowerID)
	borrower.OTPverified = true
	borrower.Save()
}
