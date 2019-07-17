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
		borrower  models.Borrower
		cekEmail  bool
		cekPhone  bool
		cekIdcard bool
		cekTaxid  bool
	)
	values := []string{}
	db := asira.App.DB
	email := c.QueryParam("email")
	phone := c.QueryParam("phone")
	idcard_number := c.QueryParam("idcard_number")
	taxid_number := c.QueryParam("taxid_number")
	cekEmail = db.Where("email = ?", email).Find(&borrower).RecordNotFound()
	cekPhone = db.Where("phone = ?", phone).Find(&borrower).RecordNotFound()
	cekIdcard = db.Where("idcard_number = ?", idcard_number).Find(&borrower).RecordNotFound()
	cekTaxid = db.Where("taxid_number = ?", taxid_number).Find(&borrower).RecordNotFound()
	// values = append(values, "Field : ")
	if cekEmail != true {
		values = append(values, "email")
	}
	if cekPhone != true {
		values = append(values, "Phone")
	}
	if cekIdcard != true {
		values = append(values, "Id Card Number")
	}
	if cekTaxid != true {
		values = append(values, "Tax Id Number")
	}

	if values != nil {
		result := "Field : " + strings.Join(values, " , ") + " Is Used"
		return returnInvalidResponse(http.StatusInternalServerError, "", result)
	}

	return c.JSON(http.StatusOK, "OK")
}
