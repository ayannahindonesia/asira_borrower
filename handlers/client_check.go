package handlers

import (
	"asira_borrower/asira"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	//Borrower payload for CheckData(...)
	Borrower struct {
		IdCardNumber string `json:"idcard_number" gorm:"column:idcard_number;type:varchar(255);unique;not null"`
		TaxIDnumber  string `json:"taxid_number" gorm:"column:taxid_number;type:varchar(255)"`
		Email        string `json:"email" gorm:"column:email;type:varchar(255);unique"`
		Phone        string `json:"phone" gorm:"column:phone;type:varchar(255);unique;not null"`
	}
)

//CheckData for unique field
func CheckData(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "CheckData"

	var (
		borrower Borrower
	)
	var values []string

	if email := c.QueryParam("email"); email != "" && !asira.App.DB.Where("email = ? AND agent_referral = 0", email).Find(&borrower).RecordNotFound() {
		values = append(values, EnglishToIndonesiaFieldsUnderscored["email"])
	}
	if phone := c.QueryParam("phone"); phone != "" && !asira.App.DB.Where("phone = ? AND agent_referral = 0", phone).Find(&borrower).RecordNotFound() {
		values = append(values, EnglishToIndonesiaFieldsUnderscored["phone"])
	}
	if idCardNumber := c.QueryParam("idcard_number"); idCardNumber != "" && !asira.App.DB.Where("idcard_number = ? AND agent_referral = 0", idCardNumber).Find(&borrower).RecordNotFound() {
		values = append(values, EnglishToIndonesiaFieldsUnderscored["idcard_number"])
	}
	if taxIDNumber := c.QueryParam("taxid_number"); taxIDNumber != "" && !asira.App.DB.Where("taxid_number = ? AND agent_referral = 0", taxIDNumber).Find(&borrower).RecordNotFound() {
		values = append(values, EnglishToIndonesiaFieldsUnderscored["taxid_number"])
	}
	if len(values) < 1 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  true,
			"message": "Ok",
		})
	}

	NLog("error", LogTag, fmt.Sprintf("fields already exist : %v", values), c.Get("user").(*jwt.Token), "", true, "")

	result := "Field : " + strings.Join(values, " , ") + " Telah Digunakan"
	return returnInvalidResponse(http.StatusUnprocessableEntity, "", result)

}
