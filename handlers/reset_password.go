package handlers

import (
	"asira/models"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func ClientResetPassword(c echo.Context) error {

	defer c.Request().Body.Close()
	borrower := models.Borrower{}
	r := c.Request()
	baseURL := c.Scheme() + "://" + r.Host
	payloadRules := govalidator.MapData{
		"email": []string{"email", "required", "unique:borrowers,email"},
	}
	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		tokenrole := "borrower"
		token, err := createJwtToken(strconv.FormatUint(borrower.ID, 10), tokenrole)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "error creating token")
		}
		to := []string{borrower.Email}
		subject := "[NO REPLY] - Reset Password Aplikasi Mobile ASIRA"
		link := baseURL + "?q=" + token
		message := "Hai Nasabah,\nIni adalah email untuk melakukan reset login akun anda.\nSilahkan klik link di bawah ini agar dapat melakukan reset login akun.\nLink ini hanya valid dalam waktu 24 jam.\n" + link + " \n Ayannah Solusi Nusantara Team"

		err = sendMail(to, subject, message)
		if err != nil {
			log.Fatal(err.Error())
		}
		return returnInvalidResponse(http.StatusOK, "", "OK")
	}
	return returnInvalidResponse(http.StatusUnprocessableEntity, "Email Not Found", "Email Not Found")
}
