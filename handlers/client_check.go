package handlers

import (
	"asira/asira"
	"asira/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
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

func CheckEmail(c echo.Context) error {
	defer c.Request().Body.Close()
	var (
		borrower models.Borrower
		token    string
		err      error
	)
	payloadRules := govalidator.MapData{
		"email": []string{"email", "required", "unique:borrowers,email"},
	}

	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		tokenrole := "borrower"
		token, err = createJwtToken(strconv.FormatUint(borrower.ID, 10), tokenrole)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "error creating token")
		}

		to := []string{borrower.Email}
		subject := "[NO REPLY] - Reset Password Aplikasi Mobile ASIRA"
		link := "https://asira.ayannnah.com/" + token
		message := "Hai Nasabah,\nIni adalah email untuk melakukan reset login akun anda.\nSilahkan klik link di bawah ini agar dapat melakukan reset login akun.\nLink ini hanya valid dalam waktu 24 jam.\n" + link + " \n Ayannah Solusi Nusantara Team"

		erro := sendMail(to, subject, message)
		if erro != nil {
			log.Fatal(erro.Error())
		}

		log.Println("Mail sent!")
		return c.JSON(http.StatusOK, borrower.Email)
	}
	log.Println(borrower)
	return returnInvalidResponse(http.StatusUnprocessableEntity, "Email Not Found", "Email Not Found")
}
