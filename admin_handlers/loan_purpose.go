package admin_handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func LoanPurposeList(c echo.Context) error {
	defer c.Request().Body.Close()

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	name := c.QueryParam("name")
	status := c.QueryParam("status")

	type Filter struct {
		Name   string `json:"name" condition:"LIKE"`
		Status string `json:"status"`
	}

	purposes := models.LoanPurpose{}
	result, err := purposes.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Name:   name,
		Status: status,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

func LoanPurposeNew(c echo.Context) error {
	defer c.Request().Body.Close()

	purpose := models.LoanPurpose{}
	payloadRules := govalidator.MapData{
		"name":   []string{"required"},
		"status": []string{"required", "loan_purpose_status"},
	}

	validate := validateRequestPayload(c, payloadRules, &purpose)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err := purpose.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat loan purpose baru")
	}

	return c.JSON(http.StatusCreated, purpose)
}

func LoanPurposeDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	loan_purpose_id, _ := strconv.ParseUint(c.Param("loan_purpose_id"), 10, 64)

	purpose := models.LoanPurpose{}
	err := purpose.FindbyID(loan_purpose_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan purpose %v tidak ditemukan", loan_purpose_id))
	}

	return c.JSON(http.StatusOK, purpose)
}

func LoanPurposePatch(c echo.Context) error {
	defer c.Request().Body.Close()

	loan_purpose_id, _ := strconv.ParseUint(c.Param("loan_purpose_id"), 10, 64)

	purpose := models.LoanPurpose{}
	err := purpose.FindbyID(loan_purpose_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan purpose %v tidak ditemukan", loan_purpose_id))
	}

	payloadRules := govalidator.MapData{
		"name":   []string{},
		"status": []string{"loan_purpose_status"},
	}

	validate := validateRequestPayload(c, payloadRules, &purpose)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = purpose.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal update loan purpose %v", loan_purpose_id))
	}

	return c.JSON(http.StatusOK, purpose)
}

func LoanPurposeDelete(c echo.Context) error {
	defer c.Request().Body.Close()

	loan_purpose_id, _ := strconv.ParseUint(c.Param("loan_purpose_id"), 10, 64)

	purpose := models.LoanPurpose{}
	err := purpose.FindbyID(loan_purpose_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan purpose %v tidak ditemukan", loan_purpose_id))
	}

	err = purpose.Delete()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, fmt.Sprintf("Gagal delete loan purpose %v", loan_purpose_id))
	}

	return c.JSON(http.StatusOK, purpose)
}
