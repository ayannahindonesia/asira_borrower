package handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

//BorrowerProfile get borrower personal profile
func AgentBorrowerProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	//get agent id
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseInt(claims["jti"].(string), 10, 64)
	borrowerID, _ := strconv.ParseUint(c.Param("borrower_id"), 10, 64)

	//cek borrower
	borrowerModel := models.Borrower{}
	err := borrowerModel.FindbyID(int(borrowerID))
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "validation error : Akun borrower agent tidak ditemukan")
	}

	//cek borrower valid, owned by agent
	if borrowerModel.AgentReferral.Int64 != agentID {
		return returnInvalidResponse(http.StatusForbidden, err, "validation error : bukan borrower agent yang valid")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}

//BorrowerProfileEdit patch data borrower personal
func AgentBorrowerProfileEdit(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseInt(claims["jti"].(string), 10, 64)
	borrowerID, _ := strconv.ParseUint(c.Param("borrower_id"), 10, 64)

	borrowerModel := models.Borrower{}
	err := borrowerModel.FindbyID(int(borrowerID))
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "validation error : Akun borrower agent tidak ditemukan")
	}

	//cek borrower valid, owned by agent
	if borrowerModel.AgentReferral.Int64 != agentID {
		return returnInvalidResponse(http.StatusForbidden, err, "validation error : bukan borrower agent yang valid")
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

	//parse payload
	validate := validateRequestPayload(c, payloadRules, &borrowerModel)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek unique for patching
	uniques := map[string]string{
		"idcard_number":      borrowerModel.IdCardNumber,
		"taxid_number":       borrowerModel.TaxIDnumber,
		"email":              borrowerModel.Email,
		"phone":              borrowerModel.Phone,
		"bank_accountnumber": borrowerModel.BankAccountNumber,
	}
	foundFields, err := checkPatchFields("borrowers", "id", borrowerModel.ID, uniques)
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error : "+foundFields)
	}

	//saving
	err = borrowerModel.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Gagal Membuat Akun")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}
