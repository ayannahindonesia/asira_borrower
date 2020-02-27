package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/middlewares"
	"asira_borrower/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	//VerifyAccountOTPrequest for requst otp
	VerifyAccountOTPrequest struct {
		Phone  string `json:"phone"`
		Try    int    `json:"try"`
		Secret string `json:"secret"`
	}

	//VerifyAccountOTPverify for verification otp
	VerifyAccountOTPverify struct {
		VerifyAccountOTPrequest
		OTPcode string `json:"otp_code"`
	}
)

//SaltOTPs prime number
var SaltOTPs = []string{"78007", "36571", "577177"}

//RegisterBorrower register borrower personal
func RegisterBorrower(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "RegisterBorrower"

	type (
		Register struct {
			Fullname string `json:"fullname"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Password string `json:"password"`
			OTPCode  string `json:"otp_code"`
		}
	)
	register := Register{}
	payloadRules := govalidator.MapData{
		"fullname": []string{"required"},
		"email":    []string{"required"},
		"phone":    []string{"required", "id_phonenumber"},
		"password": []string{"required"},
		"otp_code": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &register)
	if validate != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "empty validate borrower",
			NLOGERR: validate}, c.Get("user").(*jwt.Token), "", false, "borrower")

		NLog("warning", LogTag, fmt.Sprintf("error validate borrower : %v", validate), c.Get("user").(*jwt.Token), "", true, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//search already exist Borrower registered by agent
	err = isBorrowerAlreadyRegistered(register.Email, register.Phone)
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "empty validate borrower",
			NLOGERR: err}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Borrower personal sudah terdaftar sebelumnya")
	}

	borrower := models.Borrower{}
	//marshalling data
	r, err := json.Marshal(register)
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "error parse borrower data",
			NLOGERR: err}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	//create 1 user for 1 borrower (nasabah personal)
	user := models.User{}

	//create new personal borrower
	json.Unmarshal(r, &borrower)

	//must set nol
	borrower.Bank = sql.NullInt64{
		Int64: 0,
		Valid: true,
	}
	borrower.AgentReferral = sql.NullInt64{
		Int64: 0,
		Valid: true,
	}

	// bypass otp
	if tryValidateOTP(SaltOTPs, register.Phone, register.OTPCode) || (asira.App.ENV == "development" && register.OTPCode == "888999") {
		borrower.OTPverified = true
	} else {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: "invalid OTP register borrower ",
			NLOGERR: register}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Invalid kode OTP")
	}

	err = borrower.Create()
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error create borrower",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	err = middlewares.SubmitKafkaPayload(borrower, "borrower_create")
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: "error submit kafka create borrower",
			NLOGERR: err}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Sinkronisasi Borrower Baru Gagal")
	}

	//save borrower_id to user entity and storing
	user.Borrower = borrower.ID
	user.Password = register.Password
	user.Create()
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error create user for borrower",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	return c.JSON(http.StatusCreated, borrower)
}

func tryValidateOTP(salts []string, phone string, comparedOTP string) bool {
	//test salted OTP
	for _, salt := range salts {
		catenate := salt + phone[len(phone)-6:]
		counter, _ := strconv.Atoi(catenate)
		if asira.App.OTP.HOTP.Verify(comparedOTP, counter) {
			return true
		}
	}
	return false
}

// RequestOTPverifyAccount request OTP from /client/
func RequestOTPverifyAccount(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "RequestOTPverifyAccount"

	otpRequest := VerifyAccountOTPrequest{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	fmt.Println(borrowerID)
	payloadRules := govalidator.MapData{
		"phone":  []string{"required", "id_phonenumber"},
		"try":    []string{"required"},
		"secret": []string{"required"},
	}

	//parse payload
	validate := validateRequestPayload(c, payloadRules, &otpRequest)
	if validate != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:         "error validate borrower",
			NLOGERR:         validate,
			NLOGQUERY:       asira.App.DB.QueryExpr(),
			"borrower_user": otpRequest}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek otp request try
	// Try, _ := strconv.Atoi(otpRequest.Try)
	Try := otpRequest.Try
	if Try < 1 || Try > len(SaltOTPs) {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "invalid Try parameter",
			NLOGERR: "try < 1 or try not valid"}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "invalid Try")
	}
	//get OTP secret
	Secret := asira.App.Config.GetString(fmt.Sprintf("%s.messaging.otp_secret", asira.App.ENV))
	if otpRequest.Secret != Secret {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "invalid Try parameter",
			NLOGERR: otpRequest}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "invalid Secret")
	}

	// combine borrower id with last 6 digit of phone as counter
	catenate := SaltOTPs[Try-1] + otpRequest.Phone[len(otpRequest.Phone)-6:]
	counter, _ := strconv.Atoi(catenate)
	otpCode := asira.App.OTP.HOTP.At(int(counter))

	//send SMS OTP
	message := fmt.Sprintf("Code OTP Registrasi anda adalah %s", otpCode)
	err := asira.App.Messaging.SendSMS(otpRequest.Phone, message)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error failed Sending SMS OTP",
			NLOGERR:   err,
			"payload": otpRequest}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending otp")
	}

	NLog("info", LogTag, map[string]interface{}{
		NLOGMSG:   "success send OTP",
		"payload": otpRequest}, c.Get("user").(*jwt.Token), "", false, "borrower")

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Terkirim"})
}

// func getRandomOpt() string {
// 	len := len(SaltOTPs)
// 	i := uint32(0)
// 	if len > 0 {
// 		i = getRandomUint32() % uint32(len)
// 	}
// 	return SaltOTPs[i]
// }

// func getRandomUint32() uint32 {
// 	x := time.Now().UnixNano()
// 	test := uint32((x >> 32) ^ x)
// 	fmt.Println("Random : ", test)
// 	return test
// }

func updateAccountOTPstatus(c echo.Context, borrowerID uint64) error {
	modelBorrower := models.Borrower{}

	LogTag := "updateAccountOTPstatus"

	//get data and check status OTPverified
	_ = modelBorrower.FindbyID(borrowerID)
	if modelBorrower.OTPverified == true {
		NLog("error", LogTag, fmt.Sprintf("error already verified borrower : %v", modelBorrower), c.Get("user").(*jwt.Token), "", false, "agent")

		return errors.New("Nasabah sudah terverifikasi")
	}

	//change status verified borrower
	modelBorrower.OTPverified = true
	err = middlewares.SubmitKafkaPayload(modelBorrower, "borrower_update")
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("error kafka submit update borrower : %v borrower ID : %v", err, borrowerID), c.Get("user").(*jwt.Token), "", false, "agent")

		modelBorrower.OTPverified = false
		return err
	}

	return nil
}
