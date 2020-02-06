package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/middlewares"
	"asira_borrower/models"
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
		"email":    []string{"required", "unique:borrowers,email"},
		"phone":    []string{"required", "id_phonenumber", "unique:borrowers,phone"},
		"password": []string{"required"},
		"otp_code": []string{"required"},
		// "nickname":              []string{},
		//"birthday": []string{"date"},
		// "gender":                []string{},
		// "idcard_number":         []string{},
		// "taxid_number":          []string{},
		// "nationality":           []string{},
		// "birthplace":            []string{},
		// "last_education":        []string{},
		// "mother_name":           []string{},
		// "marriage_status":       []string{},
		// "spouse_name":           []string{},
		// "spouse_birthday":       []string{},
		// "spouse_lasteducation":  []string{},
		// "dependants":            []string{},
		// "address":               []string{},
		// "province":              []string{},
		// "city":                  []string{},
		// "neighbour_association": []string{},
		// "hamlets":               []string{},
		// "home_phonenumber":      []string{},
		// "subdistrict":           []string{},
		// "urban_village":         []string{},
		// "home_ownership":        []string{},
		// "lived_for":             []string{},
		// "occupation":            []string{},
		// "employee_id":           []string{},
		// "employer_name":         []string{},
		// "employer_address":      []string{},
		// "department":            []string{},
		// "been_workingfor":       []string{},
		// "direct_superiorname":   []string{},
		// "employer_number":       []string{},
		// "monthly_income":        []string{},
		// "other_income":          []string{},
		// "other_incomesource":    []string{},
		// "field_of_work":         []string{},
		// "related_personname":    []string{},
		// "related_relation":      []string{},
		// "related_phonenumber":   []string{},
		// "related_homenumber":    []string{},
		// "bank":                  []string{},
		// "bank_accountnumber":    []string{},
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

	// //upload image id card
	// IdCardImage, err := uploadImageS3Formatted("ktp", register.IdCardImage)
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal : IDCardImage failed to upload")
	// }

	// //upload image tax card
	// TaxIDImage, err := uploadImageS3Formatted("tax", register.TaxIDImage)
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal : TaxIDImage failed to upload")
	// }

	//create 1 user for 1 borrower (nasabah personal)
	user := models.User{}

	// //search already exist Borrower registered by agent
	// err = isBorrowerAlreadyRegistered(register.IdCardNumber)
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "Borrower personal sudah terdaftar sebelumnya")
	// }

	// //check manual fields if not unique
	// var fields = map[string]string{
	// 	"phone": register.Phone,
	// 	"email": register.Email,
	// }
	// fieldsFound, err := checkUniqueFields(register.IdCardNumber, fields)
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "data sudah ada sebelumnya : "+fieldsFound)
	// }

	//TODO: cek OTP

	//create new personal borrower
	json.Unmarshal(r, &borrower)
	// encryptPassphrase := asira.App.Config.GetString(fmt.Sprintf("%s.passphrase", asira.App.ENV))
	// borrower.IdCardImage, err = encrypt(IdCardImage, encryptPassphrase)
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi Id card gagal")
	// }
	// borrower.TaxIDImage, err = encrypt(TaxIDImage, encryptPassphrase)
	// if err != nil {
	// 	return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi NPWP gagal")
	// }

	// borrower.Bank = sql.NullInt64{
	// 	Int64: int64(register.Bank),
	// 	Valid: true,
	// }
	// borrower.AgentReferral = sql.NullInt64{
	// 	Int64: 0,
	// 	Valid: true,
	// }
	err = borrower.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}
	err = middlewares.SubmitKafkaPayload(borrower, "borrower_create")
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Sinkronisasi Borrower Baru Gagal")
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
		err = updateAccountOTPstatus(borrowerID)
		if err != nil {
			return returnInvalidResponse(http.StatusBadRequest, err, "gagal mengubah otp borrower")
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Verified"})
	}

	// bypass otp
	if asira.App.ENV == "development" && otpVerify.OTPcode == "888999" {
		err = updateAccountOTPstatus(borrowerID)
		if err != nil {
			return returnInvalidResponse(http.StatusBadRequest, err, "gagal mengubah otp borrower")
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Verified"})
	}

	return returnInvalidResponse(http.StatusBadRequest, "", "OTP salah")
}

func updateAccountOTPstatus(borrowerID uint64) error {
	modelBorrower := models.Borrower{}

	//get data and check status OTPverified
	_ = modelBorrower.FindbyID(borrowerID)
	if modelBorrower.OTPverified == true {
		return errors.New("Nasabah sudah terverifikasi")
	}

	//change status verified borrower
	modelBorrower.OTPverified = true
	err = middlewares.SubmitKafkaPayload(modelBorrower, "borrower_update")
	if err != nil {
		modelBorrower.OTPverified = false
		return err
	}

	return nil
}
