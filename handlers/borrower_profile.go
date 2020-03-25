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
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "error get borrower profile",
			NLOGERR:       err,
			NLOGQUERY:     asira.App.DB.QueryExpr(),
			"borrower_id": borrowerID}, c.Get("user").(*jwt.Token), "", false, "borrower")

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
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "not valid borrower personal",
			NLOGERR:       err,
			NLOGQUERY:     asira.App.DB.QueryExpr(),
			"borrower_id": borrowerID}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}
	origin := borrowerModel

	payloadRules := govalidator.MapData{
		"fullname":              []string{"required"},
		"image_profile":         []string{},
		"nickname":              []string{},
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

		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "validation error",
			NLOGERR:   validate,
			"payload": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek bank
	if !validBankID(borrowerModel.Bank.Int64) {

		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "validation error",
			NLOGERR:   "invalid bank ID",
			"payload": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "bank tidak valid")
	}

	//cek unique for patching
	var fields = map[string]string{
		"phone":              borrowerModel.Phone,
		"email":              borrowerModel.Email,
		"taxid_number":       borrowerModel.TaxIDnumber,
		"bank_accountnumber": borrowerModel.BankAccountNumber,
		"idcard_number":      borrowerModel.IdCardNumber,
	}
	//custom patch, coz personal and agent's might be exist
	fieldsFound, err := checkFieldsBorrowersPersonal(borrowerModel.ID, fields)
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:        "error validate patching borrower",
			NLOGERR:        err,
			"fields-found": fieldsFound,
			"borrower":     borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "data sudah ada sebelumnya : "+fieldsFound)
	}

	//get passphrase encryption
	encryptPassphrase := asira.App.Config.GetString(fmt.Sprintf("%s.passphrase", asira.App.ENV))

	//upload image id card
	if borrowerModel.IdCardImage != "" || len(borrowerModel.IdCardImage) != 0 {
		IDCardImage, err := uploadImageS3Formatted("ktp", borrowerModel.IdCardImage)
		if err != nil {
			NLog("error", LogTag, map[string]interface{}{
				NLOGMSG:    "error uploading ID Card image",
				NLOGERR:    err,
				"borrower": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gambar KTP gagal diunggah")
		}

		//encrypt image url
		borrowerModel.IdCardImage, err = encrypt(IDCardImage, encryptPassphrase)
		if err != nil {
			NLog("error", LogTag, map[string]interface{}{
				NLOGMSG:    "error encrypting ID Card image",
				NLOGERR:    err,
				"borrower": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi KTP gagal")
		}
	}

	//upload image tax card
	if borrowerModel.TaxIDImage != "" || len(borrowerModel.TaxIDImage) != 0 {
		TaxIDImage, err := uploadImageS3Formatted("tax", borrowerModel.TaxIDImage)
		if err != nil {
			NLog("error", LogTag, map[string]interface{}{
				NLOGMSG:    "error uploading Tax ID image",
				NLOGERR:    err,
				"borrower": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gambar NPWP gagal diunggah")
		}

		//encrypt image url
		borrowerModel.TaxIDImage, err = encrypt(TaxIDImage, encryptPassphrase)
		if err != nil {
			NLog("error", LogTag, map[string]interface{}{
				NLOGMSG:    "error encrypting Tax ID image",
				NLOGERR:    err,
				"borrower": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi NPWP gagal")
		}
	}

	if borrowerModel.ImageProfile != "" || len(borrowerModel.ImageProfile) != 0 {
		//upload ImageProfile from base64 and get response lintasarta url
		borrowerModel.ImageProfile, err = uploadImageS3Formatted("borrprofile", borrowerModel.ImageProfile)
		if err != nil {
			NLog("error", LogTag, map[string]interface{}{
				NLOGMSG:            "error uploading Image Profile",
				NLOGERR:            err,
				"borrower_payload": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "agent")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Ubah Data Borrower Gagal : Image profil failed to upload")
		}
	}

	// err = borrowerModel.Save()
	err = middlewares.SubmitKafkaPayload(borrowerModel, "borrower_update")
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:    "error submitting to kafka after creating borrower",
			NLOGERR:    err,
			"borrower": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal update Borrower")
	}

	NLog("info", LogTag, map[string]interface{}{
		NLOGMSG:    "succcess borrower edit profile",
		"borrower": borrowerModel}, c.Get("user").(*jwt.Token), "", false, "borrower")
	NAudittrail(origin, borrowerModel, token, "borrower", fmt.Sprint(borrowerModel.ID), "borrower edit profile", "borrower")

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
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "not valid borrower personal",
			NLOGERR:       err,
			"borrower_id": borrowerID}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun bukan borrower personal")
	}

	payloadRules := govalidator.MapData{
		"password": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &userBorrower)
	if validate != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "error validation",
			NLOGERR: validate}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(userBorrower.Password), bcrypt.DefaultCost)
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "error generate passoword",
			NLOGERR: err}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return err
	}

	//update to new password
	userBorrower.Password = string(passwordByte)
	err = userBorrower.Save()
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:         "Failed update borrower password",
			NLOGERR:         err,
			NLOGQUERY:       asira.App.DB.QueryExpr(),
			"borrower_user": userBorrower}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Ubah Password Gagal")
	}

	NAudittrail(models.Loan{}, userBorrower, token, "borrower", fmt.Sprint(userBorrower.ID), "borrower change password", "borrower")

	responseBody := map[string]interface{}{
		"status":  true,
		"message": "Ubah Passord berhasil",
	}
	return c.JSON(http.StatusOK, responseBody)
}
