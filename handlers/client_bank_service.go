package handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//ClientBankServices get service
func ClientBankServices(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "ClientBankServices"

	bankService := models.Service{}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	// filters
	type Filter struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	result, err := bankService.PagedFindFilter(page, rows, orderby, sort, &Filter{
		Name:   c.QueryParam("name"),
		Status: c.QueryParam("status"),
	})

	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("query result error : %v", err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}

//ClientBankServicebyID get service detail
func ClientBankServicebyID(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "ClientBankServicebyID"

	bankService := models.Service{}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := bankService.FindbyID(id)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("service not found : %v", err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, err, "service tidak ditemukan")
	}

	return c.JSON(http.StatusOK, bankService)
}
