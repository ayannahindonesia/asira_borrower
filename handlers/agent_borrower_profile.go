package handlers

import (
	"asira_borrower/middlewares"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

//AgentBorrowerProfile get borrower personal profile
func AgentBorrowerProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentBorrowerProfile"

	//get agent id
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseInt(claims["jti"].(string), 10, 64)
	borrowerID, _ := strconv.ParseUint(c.Param("borrower_id"), 10, 64)

	//cek borrower
	borrowerModel := models.Borrower{}
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("not valid borrower : %v borrower ID : %v", err, borrowerID), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusNotFound, err, "validation error : Akun borrower agent tidak ditemukan")
	}

	//cek borrower valid, owned by agent
	if borrowerModel.AgentReferral.Int64 != agentID {
		NLog("error", LogTag, fmt.Sprintf("borrower %v not owned by agent ID : %v", borrowerID, agentID), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusForbidden, err, "validation error : bukan borrower agent yang valid")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}

//AgentBorrowerProfileEdit patch data borrower personal
func AgentBorrowerProfileEdit(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentBorrowerProfileEdit"

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseInt(claims["jti"].(string), 10, 64)
	borrowerID, _ := strconv.ParseUint(c.Param("borrower_id"), 10, 64)

	borrowerModel := models.Borrower{}
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		NLog("warning", LogTag, fmt.Sprintf("not valid borrower : %v borrower ID : %v", err, borrowerID), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusNotFound, err, "validation error : Akun borrower agent tidak ditemukan")
	}

	//cek borrower valid, owned by agent
	if borrowerModel.AgentReferral.Int64 != agentID {
		NLog("error", LogTag, fmt.Sprintf("borrower %v not owned by agent ID : %v", borrowerID, agentID), c.Get("user").(*jwt.Token), "", false, "agent")

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
	var borrowerPayload models.Borrower
	validate := validateRequestPayload(c, payloadRules, &borrowerPayload)
	if validate != nil {
		NLog("warning", LogTag, fmt.Sprintf("error validation : %v", validate), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek unique for patching
	uniques := map[string]string{
		"taxid_number":       borrowerPayload.TaxIDnumber,
		"email":              borrowerPayload.Email,
		"phone":              borrowerPayload.Phone,
		"bank_accountnumber": borrowerPayload.BankAccountNumber,
	}
	foundFields, err := checkUniqueFields(borrowerModel.IdCardNumber, uniques)
	if err != nil {
		NLog("warning", LogTag, fmt.Sprintf("data already exist : %v", foundFields), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error : "+foundFields)
	}

	if borrowerPayload.MonthlyIncome > 0 {
		borrowerModel.MonthlyIncome = borrowerPayload.MonthlyIncome

	}
	if borrowerPayload.OtherIncome > 0 {
		borrowerModel.OtherIncome = borrowerPayload.OtherIncome

	}
	if len(borrowerPayload.OtherIncomeSource) > 0 {
		borrowerModel.OtherIncomeSource = borrowerPayload.OtherIncomeSource

	}
	//saving
	// err = borrowerModel.Save()
	err = middlewares.SubmitKafkaPayload(borrowerModel, "borrower_update")
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("error kafka submit update borrower : %v borrower ID : %v", err, borrowerID), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal update Borrower")
	}

	return c.JSON(http.StatusOK, borrowerModel)
}
