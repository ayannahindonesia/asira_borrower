package handlers

import (
	"asira/models"
	"net/http"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func RegisterBorrower(c echo.Context) error {
	defer c.Request().Body.Close()

	borrower := models.Borrower{}

	payloadRules := govalidator.MapData{
		"fullname":              []string{"required"},
		"gender":                []string{"required"},
		"idcard_number":         []string{"required"},
		"taxid_number":          []string{},
		"email":                 []string{"email"},
		"birthday":              []string{"date"},
		"birthplace":            []string{"required"},
		"last_education":        []string{"required"},
		"mother_name":           []string{"required"},
		"phone":                 []string{"required"},
		"marriage_status":       []string{"required"},
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
		"home_ownership":        []string{"required"},
		"lived_for":             []string{"required"},
		"occupation":            []string{"required"},
		"employee_id":           []string{},
		"employer_name":         []string{"required"},
		"employer_address":      []string{"required"},
		"department":            []string{"required"},
		"been_workingfor":       []string{"required"},
		"direct_superiorname":   []string{},
		"employer_number":       []string{"required"},
		"monthly_income":        []string{"required"},
		"other_income":          []string{},
		"other_incomesource":    []string{},
		"field_of_work":         []string{"required"},
		"related_personname":    []string{"required"},
		"related_relation":      []string{"required"},
		"related_phonenumber":   []string{"required"},
		"related_homenumber":    []string{},
		"password":              []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	newBorrower, err := borrower.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "create new borrower failed")
	}

	return c.JSON(http.StatusCreated, newBorrower)
}
