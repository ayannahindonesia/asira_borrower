package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/middlewares"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

//BorrowerPersonalResponse custom response
type BorrowerPersonalResponse struct {
	models.Borrower
	LoanStatus string `json:"loan_status"`
}

//BorrowerProfile get borrower personal profile
func BorrowerProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "BorrowerProfile"

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	//check current borrower
	borrower := BorrowerPersonalResponse{}
	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	//manual query
	db := asira.App.DB

	//query loan from borrowers
	LoanStatusQuery := fmt.Sprintf("CASE WHEN (SELECT COUNT(id) FROM loans l WHERE l.borrower = borrowers.id AND status IN ('%s', '%s') AND (due_date IS NULL OR due_date = '0001-01-01 00:00:00+00' OR NOW() < l.due_date + make_interval(days => 1))) > 0  THEN '%s' ELSE '%s' END", "approved", "processing", "active", "inactive")

	//gen query
	db = db.Table("borrowers").
		Select("borrowers.*, "+LoanStatusQuery+" as loan_status").
		Where("borrowers.id = ?", borrowerID)

	err = db.Find(&borrower).Error
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("error get borrower profile : %v borrower id : %v", err, borrowerID), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}
	return c.JSON(http.StatusOK, borrower)
}

//BorrowerProfileEdit patch data borrower personal
func BorrowerProfileEdit(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "BorrowerProfileEdit"

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		NLog("warning", LogTag, fmt.Sprintf("not valid borrower personal : %v borrower id : %v", err, borrowerID), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}
	origin := borrowerModel

	payloadRules := govalidator.MapData{
		"fullname":              []string{},
		"gender":                []string{},
		"idcard_number":         []string{"required"},
		"taxid_number":          []string{},
		"email":                 []string{"email"},
		"birthday":              []string{"date"},
		"birthplace":            []string{},
		"last_education":        []string{},
		"mother_name":           []string{"required"},
		"phone":                 []string{},
		"marriage_status":       []string{},
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
		"home_ownership":        []string{},
		"lived_for":             []string{},
		"occupation":            []string{"required"},
		"employee_id":           []string{},
		"employer_name":         []string{},
		"employer_address":      []string{},
		"department":            []string{},
		"been_workingfor":       []string{},
		"direct_superiorname":   []string{},
		"employer_number":       []string{},
		"monthly_income":        []string{},
		"other_income":          []string{},
		"other_incomesource":    []string{},
		"field_of_work":         []string{},
		"related_personname":    []string{},
		"related_relation":      []string{},
		"related_phonenumber":   []string{},
		"related_homenumber":    []string{},
		"related_address":       []string{},
		"bank":                  []string{},
		"bank_accountnumber":    []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &borrowerModel)
	if validate != nil {
		NLog("warning", LogTag, fmt.Sprintf("validation error : %v", validate), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek unique for patching
	var fields = map[string]string{
		"phone":              borrowerModel.Phone,
		"email":              borrowerModel.Email,
		"taxid_number":       borrowerModel.TaxIDnumber,
		"bank_accountnumber": borrowerModel.BankAccountNumber,
	}
	//custom patch, coz personal and agent's might be exist
	fieldsFound, err := checkPatchFieldsBorrowers(borrowerModel.ID, borrowerModel.IdCardNumber, fields)
	if err != nil {
		NLog("warning", LogTag, fmt.Sprintf("error validate patching borrower : %v (%v) borrower : %v", err, fieldsFound, borrowerModel), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "data sudah ada sebelumnya : "+fieldsFound)
	}

	//get passphrase encryption
	encryptPassphrase := asira.App.Config.GetString(fmt.Sprintf("%s.passphrase", asira.App.ENV))

	//upload image id card
	if borrowerModel.IdCardImage != "" || len(borrowerModel.IdCardImage) != 0 {
		IDCardImage, err := uploadImageS3Formatted("ktp", borrowerModel.IdCardImage)
		if err != nil {
			NLog("error", LogTag, fmt.Sprintf("error uploading ID Card image : %v borrower : %v", err, borrowerModel), c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gambar KTP gagal diunggah")
		}

		//encrypt image url
		borrowerModel.IdCardImage, err = encrypt(IDCardImage, encryptPassphrase)
		if err != nil {
			NLog("error", LogTag, fmt.Sprintf("error encrypting ID Card image : %v borrower : %v", err, borrowerModel), c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi KTP gagal")
		}
	}

	//upload image tax card
	if borrowerModel.TaxIDImage != "" || len(borrowerModel.TaxIDImage) != 0 {
		TaxIDImage, err := uploadImageS3Formatted("tax", borrowerModel.TaxIDImage)
		if err != nil {
			NLog("error", LogTag, fmt.Sprintf("error uploading Tax ID image : %v borrower : %v", err, borrowerModel), c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gambar NPWP gagal diunggah")
		}

		//encrypt image url
		borrowerModel.TaxIDImage, err = encrypt(TaxIDImage, encryptPassphrase)
		if err != nil {
			NLog("error", LogTag, fmt.Sprintf("error encrypting Tax ID image : %v borrower : %v", err, borrowerModel), c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi NPWP gagal")
		}
	}

	// err = borrowerModel.Save()
	err = middlewares.SubmitKafkaPayload(borrowerModel, "borrower_update")
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("error submitting to kafka after creating borrower : %v borrower : %v", err, borrowerModel), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal update Borrower")
	}

	NAudittrail(origin, borrowerModel, token, "borrower", fmt.Sprint(borrowerModel.ID), "borrower edit profile")

	NLog("event", LogTag, fmt.Sprintf("borrower edit profile : %v", borrowerModel), c.Get("user").(*jwt.Token), "", false, "borrower")

	return c.JSON(http.StatusOK, borrowerModel)
}

//BorrowerChangePassword update borrower personal password
func BorrowerChangePassword(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "BorrowerChangePassword"

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	//get password from users entity/table
	userBorrower := models.User{}
	err = userBorrower.FindbyBorrowerID(borrowerID)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("not valid borrower personal : %v borrower id : %v", err, borrowerID), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun bukan borrower personal")
	}

	payloadRules := govalidator.MapData{
		"password": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &userBorrower)
	if validate != nil {
		NLog("warning", LogTag, fmt.Sprintf("error validation : %v", validate), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(userBorrower.Password), bcrypt.DefaultCost)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("error generate passoword : %v", err), c.Get("user").(*jwt.Token), "", false, "borrower")

		return err
	}

	//update to new password
	userBorrower.Password = string(passwordByte)
	err = userBorrower.Save()
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("Failed update borrower password : %v ; borrower : %v", err, userBorrower), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Ubah Password Gagal")
	}

	NAudittrail(models.Loan{}, userBorrower, token, "borrower", fmt.Sprint(userBorrower.ID), "borrower change password")

	responseBody := map[string]interface{}{
		"status":  true,
		"message": "Ubah Passord berhasil",
	}
	return c.JSON(http.StatusOK, responseBody)
}
