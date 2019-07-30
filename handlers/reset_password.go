package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
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
	baseURL := c.Scheme() + "://" + r.Host + "/"
	payloadRules := govalidator.MapData{
		"email": []string{"email", "unique:borrowers,email"},
	}
	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		tokenrole := "borrower"
		token, err := createJwtToken(strconv.FormatUint(borrower.ID, 10), tokenrole)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "error creating token")
		}
		if borrower.Email == "" {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Email Not Found")
		}
		to := borrower.Email
		subject := "[NO REPLY] - Reset Password Aplikasi Mobile ASIRA"
		link := baseURL + "?q=" + token
		message := "Hai Nasabah,\nIni adalah email untuk melakukan reset login akun anda.\nSilahkan klik link di bawah ini agar dapat melakukan reset login akun.\nLink ini hanya valid dalam waktu 24 jam.\n" + link + " \n Ayannah Solusi Nusantara Team"

		Config := asira.App.Config.GetStringMap(fmt.Sprintf("%s.mailer", asira.App.ENV))
		log.Println(Config)
		err = sendMail(to, subject, message)
		if err != nil {
			log.Println(err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Link reset password telah dikirim ke email", "status": true})
	}
	return returnInvalidResponse(http.StatusNotFound, "", "Email Tidak Ditemukan")
}
