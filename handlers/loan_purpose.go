package handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//LoanPurposeList get loan purpose list
func LoanPurposeList(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "LoanPurposeList"

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
		NLog("warning", LogTag, fmt.Sprintf("error get Loan Purpose list  : %v", err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, err, "pencarian tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}

//LoanPurposeDetail get loan purpose detail
func LoanPurposeDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "LoanPurposeDetail"

	loanPurposeID, _ := strconv.ParseUint(c.Param("loan_purpose_id"), 10, 64)

	purpose := models.LoanPurpose{}
	err := purpose.FindbyID(loanPurposeID)
	if err != nil {
		NLog("warning", LogTag, fmt.Sprintf("loan purpose %v not found : %v", loanPurposeID, err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan purpose %v tidak ditemukan", loanPurposeID))
	}

	return c.JSON(http.StatusOK, purpose)
}
