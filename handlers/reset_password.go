package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"
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
		asira.App.DB.Where("email = ?", borrower.Email).Find(&borrower)

		uuid := models.Uuid_Reset_Password{
			Borrower: sql.NullInt64{
				Int64: int64(borrower.BaseModel.ID),
				Valid: true,
			},
		}
		Reset, err := uuid.Create()
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
		}
		if borrower.OTPverified != true {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Akun Anda Belum di Verifikasi, Silahkan Verifikasi Terlebih Dahulu")
		}
		if borrower.Email == "" {
			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Email Tidak Ditemukan atau Email anda belum terdaftar")
		}
		to := borrower.Email
		subject := "[NO REPLY] - Reset Password Aplikasi ASIRA"
		link := baseURL + "/deepLinks/" + Reset.UUID
		message := "Hai Nasabah,\n\nIni adalah email untuk melakukan reset login akun anda. Silahkan klik link di bawah ini agar dapat melakukan reset login akun.\nLink ini hanya valid dalam waktu 24 jam.\n" + link + " \n\n\n Ayannah Solusi Nusantara Team"

		err = sendMail(to, subject, message)
		if err != nil {
			log.Println(err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Link reset password telah dikirim ke email, silahkan cek email anda"})
	}
	return returnInvalidResponse(http.StatusNotFound, "", "Email Tidak Ditemukan atau Email anda belum terdaftar")
}

func ChangePassword(c echo.Context) error {
	defer c.Request().Body.Close()
	type (
		Reset struct {
			Password string `json:"password"`
			UUID     string `json:"uuid"`
		}
	)
	reset := Reset{}
	validate := validateRequestPayload(c, payloadRules, &reset)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}
	//check UUID
	uuid_reset_password := models.Uuid_Reset_Password{}
	type Filter struct {
		UUID string `json:"uuid"`
	}
	result, err := uuid_reset_password.FilterSearchSingle(&Filter{
		UUID: reset.UUID,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("UUID : %v tidak ditemukan", reset.UUID))
	}
	//check Borrower ID
	borrowerModel := models.Borrower{}
	borrower, err := borrowerModel.FindbyID(result.Borrower)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun Tidak ditemukan")
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(borrower.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	borrower.Password = string(passwordByte)
	_, err = borrower.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Ubah Password Gagal")
	}
	responseBody := map[string]interface{}{
		"status":  true,
		"message": "Ubah Passord berhasil",
	}
	return c.JSON(http.StatusOK, responseBody)
}
