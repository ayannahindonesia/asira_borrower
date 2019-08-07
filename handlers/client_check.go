package handlers

import (
	"asira_borrower/asira"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

type (
	Borrower struct {
		IdCardNumber string `json:"idcard_number" gorm:"column:idcard_number;type:varchar(255);unique;not null"`
		TaxIDnumber  string `json:"taxid_number" gorm:"column:taxid_number;type:varchar(255)"`
		Email        string `json:"email" gorm:"column:email;type:varchar(255);unique"`
		Phone        string `json:"phone" gorm:"column:phone;type:varchar(255);unique;not null"`
	}
)

func CheckData(c echo.Context) error {
	defer c.Request().Body.Close()
	var (
		borrower Borrower
	)
	var values []string

	if email := c.QueryParam("email"); email != "" && !asira.App.DB.Where("email = ?", email).Find(&borrower).RecordNotFound() {
		values = append(values, "Email")
	}
	if phone := c.QueryParam("phone"); phone != "" && !asira.App.DB.Where("phone = ?", phone).Find(&borrower).RecordNotFound() {
		values = append(values, "Phone")
	}
	if idcard_number := c.QueryParam("idcard_number"); idcard_number != "" && !asira.App.DB.Where("idcard_number = ?", idcard_number).Find(&borrower).RecordNotFound() {
		values = append(values, "Id Card Number")
	}
	if taxid_number := c.QueryParam("taxid_number"); taxid_number != "" && !asira.App.DB.Where("taxid_number = ?", taxid_number).Find(&borrower).RecordNotFound() {
		values = append(values, "Tax Id Number")
	}
	if len(values) < 1 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  true,
			"message": "Ok",
		})
	}
	result := "Field : " + strings.Join(values, " , ") + " Telah Digunakan"
	return returnInvalidResponse(http.StatusUnprocessableEntity, "", result)

}
