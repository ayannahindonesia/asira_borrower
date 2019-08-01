package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
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
		"email": []string{"email", "unique:borrowers,email"},
	}
	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		data := asira.App.DB.Where("email = ?", borrower.Email).Find(&borrower)
		log.Println(data)
		tokenrole := "borrower"
		token, err := createJwtToken(strconv.FormatUint(borrower.ID, 10), tokenrole)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "error creating token")
		}
		if borrower.OTPverified != true {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Akun Anda Belum di Verifikasi, Silahkan Verifikasi Terlebih Dahulu")
		}
		if borrower.Email == "" {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Email Tidak Ditemukan atau Email anda belum terdaftar")
		}
		to := borrower.Email
		subject := "[NO REPLY] - Reset Password Aplikasi ASIRA"
		link := baseURL + "?q=" + token
		message := "Hai Nasabah,\n\nIni adalah email untuk melakukan reset login akun anda. Silahkan klik link di bawah ini agar dapat melakukan reset login akun.\nLink ini hanya valid dalam waktu 24 jam.\n" + link + " \n\n\n Ayannah Solusi Nusantara Team"

		err = sendMail(to, subject, message)
		if err != nil {
			log.Println(err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Link reset password telah dikirim ke email, silahkan cek email anda", "status": true})
	}
	return returnInvalidResponse(http.StatusNotFound, "", "Email Tidak Ditemukan atau Email anda belum terdaftar")
}
