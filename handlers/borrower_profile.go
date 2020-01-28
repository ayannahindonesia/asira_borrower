package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

type BorrowerPersonalResponse struct {
	models.Borrower
	LoanStatus string `json:"loan_status"`
}

//BorrowerProfile get borrower personal profile
func BorrowerProfile(c echo.Context) error {
	defer c.Request().Body.Close()

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
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}
	return c.JSON(http.StatusOK, borrower)
}

//BorrowerProfileEdit patch data borrower personal
func BorrowerProfileEdit(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	payloadRules := govalidator.MapData{
		"fullname":              []string{},
		"gender":                []string{},
		"idcard_number":         []string{},
		"taxid_number":          []string{},
		"email":                 []string{"email"},
		"birthday":              []string{"date"},
		"birthplace":            []string{},
		"last_education":        []string{},
		"mother_name":           []string{},
		"phone":                 []string{},
		"marriage_status":       []string{},
		"spouse_name":           []string{},
		"spouse_birthday":       []string{"date"},
		"spouse_lasteducation":  []string{},
		"dependants":            []string{},
		"address":               []string{},
		"province":              []string{},
		"city":                  []string{},
		"neighbour_association": []string{},
		"hamlets":               []string{},
		"home_phonenumber":      []string{},
		"subdistrict":           []string{},
		"urban_village":         []string{},
		"home_ownership":        []string{},
		"lived_for":             []string{},
		"occupation":            []string{},
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
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek unique for patching
	var fields = map[string]string{
		"phone":              borrowerModel.Phone,
		"email":              borrowerModel.Email,
		"taxid_number":       borrowerModel.TaxIDnumber,
		"bank_accountnumber": borrowerModel.BankAccountNumber,
	}
	fieldsFound, err := checkUniqueFields(borrowerModel.IdCardNumber, fields)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "data sudah ada sebelumnya : "+fieldsFound)
	}

	err = borrowerModel.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Gagal Membuat Akun")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}

//BorrowerChangePassword update borrower personal password
func BorrowerChangePassword(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	//get password from users entity/table
	userBorrower := models.User{}
	err = userBorrower.FindbyBorrowerID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun bukan borrower personal")
	}

	payloadRules := govalidator.MapData{
		"password": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &userBorrower)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(userBorrower.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//update to new password
	userBorrower.Password = string(passwordByte)
	err = userBorrower.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Ubah Password Gagal")
	}
	responseBody := map[string]interface{}{
		"status":  true,
		"message": "Ubah Passord berhasil",
	}
	return c.JSON(http.StatusOK, responseBody)
}
