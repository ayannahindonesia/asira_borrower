package handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
)

func BorrowerProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}

func BorrowerProfileEdit(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}
	password := borrowerModel.Password

	payloadRules := govalidator.MapData{
		"fullname":              []string{},
		"gender":                []string{},
		"idcard_number":         []string{"unique_edit:borrowers,idcard_number"},
		"taxid_number":          []string{"unique_edit:borrowers,taxid_number"},
		"email":                 []string{"email", "unique:borrowers,taxid_number"},
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

	borrowerModel.Password = password
	err = borrowerModel.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Gagal Membuat Akun")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}

func BorrowerChangePassword(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun Tidak ditemukan")
	}

	payloadRules := govalidator.MapData{
		"password": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &borrowerModel)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(borrowerModel.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	borrowerModel.Password = string(passwordByte)
	err = borrowerModel.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Ubah Password Gagal")
	}
	responseBody := map[string]interface{}{
		"status":  true,
		"message": "Ubah Passord berhasil",
	}
	return c.JSON(http.StatusOK, responseBody)
}
