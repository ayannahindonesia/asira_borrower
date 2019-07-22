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
	var values []string
	db := asira.App.DB

	if email := c.QueryParam("email"); email != "" && !db.Where("email = ?", email).Find(&borrower).RecordNotFound() {
		values = append(values, "Email")
	}
	if phone := c.QueryParam("phone"); phone != "" && !db.Where("phone = ?", phone).Find(&borrower).RecordNotFound() {
		values = append(values, "Phone")
	}
	if idcard_number := c.QueryParam("idcard_number"); idcard_number != "" && !db.Where("idcard_number = ?", idcard_number).Find(&borrower).RecordNotFound() {
		values = append(values, "Id Card Number")
	}
	if taxid_number := c.QueryParam("taxid_number"); taxid_number != "" && !db.Where("taxid_number = ?", taxid_number).Find(&borrower).RecordNotFound() {
		values = append(values, "Tax Id Number")
	}
	if len(values) < 1 {
		return c.JSON(http.StatusOK, "OK")
	}
	result := "Field : " + strings.Join(values, " , ") + " Is Used"
	return returnInvalidResponse(http.StatusUnprocessableEntity, "", result)

}
