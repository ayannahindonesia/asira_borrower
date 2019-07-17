package handlers

import (
	"asira/asira"
	"asira/models"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func CheckData(c echo.Context) error {
	// defer c.Request().Body.Close()
	var (
		borrower models.Borrower
	)
	values := []string{}
	db := asira.App.DB
	email := c.QueryParam("email")
	phone := c.QueryParam("phone")
	idcard_number := c.QueryParam("idcard_number")
	taxid_number := c.QueryParam("taxid_number")
	if db.Where("email = ?", email).Find(&borrower).RecordNotFound() {
		values = append(values, "Email")
	}
	if db.Where("phone = ?", phone).Find(&borrower).RecordNotFound() {
		values = append(values, "Phone")
	}
	if db.Where("idcard_number = ?", idcard_number).Find(&borrower).RecordNotFound() {
		values = append(values, "Id Card Number")
	}
	if db.Where("taxid_number = ?", taxid_number).Find(&borrower).RecordNotFound() {
		values = append(values, "Tax Id Number")
	}

	if values != nil {
		result := "Field : " + strings.Join(values, " , ") + " Is Used"
		return returnInvalidResponse(http.StatusInternalServerError, "", result)
	}

	return c.JSON(http.StatusOK, "OK")
}
