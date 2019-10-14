package admin_handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
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

func LoanPurposeDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	loan_purpose_id, _ := strconv.Atoi(c.Param("loan_purpose_id"))

	purpose := models.LoanPurpose{}
	err := purpose.FindbyID(loan_purpose_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan purpose %v tidak ditemukan", loan_purpose_id))
	}

	return c.JSON(http.StatusOK, purpose)
}
