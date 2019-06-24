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
	OTPrequest struct {
		EntityType string `json:"entity_type"`
		EntityID   uint64 `json:"entity_id"`
	}
	OTPverify struct {
		OTPrequest
		OTPcode string `json:"otp_code"`
	}
)

func RequestOTPverification(c echo.Context) error {
	defer c.Request().Body.Close()

	otpRequest := OTPrequest{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	payloadRules := govalidator.MapData{
		"entity_type": []string{"otp_entity_types"},
		"entity_id":   []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &otpRequest)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	var (
		otpCode string
		counter string
	)
	switch otpRequest.EntityType {
	default:
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	case "loan":
		otpCode = asira.App.OTP.HOTP.At(int(otpRequest.EntityID))
		counter = strconv.Itoa(int(otpRequest.EntityID))
	case "borrower":
		otpCode = asira.App.OTP.HOTP.At(borrowerID)
		counter = strconv.Itoa(borrowerID)
	}
	log.Printf("generated code for counter %s : %s", counter, otpCode)

	return c.JSON(http.StatusOK, "OTP Sent")
}

func VerifyOTP(c echo.Context) error {
	defer c.Request().Body.Close()

	otpVerify := OTPverify{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	payloadRules := govalidator.MapData{
		"entity_type": []string{"otp_entity_types,required"},
		"entity_id":   []string{},
		"otp_code":    []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &otpVerify)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	switch otpVerify.EntityType {
	case "loan":
		if asira.App.OTP.HOTP.Verify(otpVerify.OTPcode, int(otpVerify.EntityID)) {
			updateOTPstatus(otpVerify.EntityType, int(otpVerify.EntityID))
			return c.JSON(http.StatusOK, "OTP Verified")
		}
	case "borrower":
		if asira.App.OTP.HOTP.Verify(otpVerify.OTPcode, borrowerID) {
			updateOTPstatus(otpVerify.EntityType, borrowerID)
			return c.JSON(http.StatusOK, "OTP Verified")
		}
	}

	return c.JSON(http.StatusBadRequest, "Wrong OTP code")
}

func updateOTPstatus(entityType string, identifier int) {
	switch entityType {
	case "loan":
		modelLoan := models.Loan{}
		loan, _ := modelLoan.FindbyID(identifier)
		loan.OTPverified = true
		loan.Save()
	case "borrower":
		modelBorrower := models.Borrower{}
		borrower, _ := modelBorrower.FindbyID(identifier)
		borrower.OTPverified = true
		borrower.Save()
	}
}
