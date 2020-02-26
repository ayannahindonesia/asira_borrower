package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//ClientBanks get banks
func ClientBanks(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "ClientBanks"

	banks := models.Bank{}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	var filter struct{}

	result, err := banks.PagedFilterSearch(page, rows, orderby, sort, &filter)

	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "query result error",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}

//ClientBankbyID get bank detail
func ClientBankbyID(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "ClientBankbyID"

	bank := models.Bank{}
	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)
	err := bank.FindbyID(bankID)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("query result error : %v", err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, err, "bank tidak ditemukan")
	}

	return c.JSON(http.StatusOK, bank)
}
