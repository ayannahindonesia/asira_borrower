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
