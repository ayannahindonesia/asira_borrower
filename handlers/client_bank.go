package handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func ClientBanks(c echo.Context) error {
	defer c.Request().Body.Close()

	banks := models.Bank{}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	var filter struct{}

	result, err := banks.PagedFilterSearch(page, rows, orderby, sort, &filter)

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}

func ClientBankbyID(c echo.Context) error {
	defer c.Request().Body.Close()

	bank := models.Bank{}
	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)
	err := bank.FindbyID(bankID)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "bank tidak ditemukan")
	}

	return c.JSON(http.StatusOK, bank)
}
