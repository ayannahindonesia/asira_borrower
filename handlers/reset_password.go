package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"

	guuid "github.com/google/uuid"
)

func ClientResetPassword(c echo.Context) error {

	defer c.Request().Body.Close()
	borrower := models.Borrower{}
	r := c.Request()
	baseURL := c.Scheme() + "://ayannah.co.id" // + r.Host
	payloadRules := govalidator.MapData{
		"email": []string{"email", "unique:borrowers,email"},
	}
	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		asira.App.DB.Where("email = ? AND agent_referral = 0", borrower.Email).Find(&borrower)
		id := guuid.New()

		uuid := models.Uuid_Reset_Password{
			Borrower: sql.NullInt64{
				Int64: int64(borrower.ID),
				Valid: true,
			},
			UUID: id.String(),
		}
		err := uuid.Create()
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
		link := baseURL + "/deepLinks/" + uuid.UUID
		message := "Hai Nasabah,\n\nIni adalah email untuk melakukan reset login akun anda. Silahkan klik link di bawah ini agar dapat melakukan reset login akun.\nLink ini hanya valid dalam waktu 24 jam.\n" + link + " \n\n\n Ayannah Solusi Nusantara Team"

		err = SendMail(to, subject, message)
		if err != nil {
			log.Println(err.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Link reset password telah dikirim ke email, silahkan cek email anda"})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Link reset password telah dikirim ke email, silahkan cek email anda"})
}

func ChangePassword(c echo.Context) error {
	defer c.Request().Body.Close()
	now := time.Now()
	type (
		Reset struct {
			Password string `json:"password"`
			UUID     string `json:"uuid"`
		}
	)
	reset := Reset{}

	payloadRules := govalidator.MapData{
		"password": []string{},
		"uuid":     []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &reset)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}
	//check UUID
	uuid_reset_password := models.Uuid_Reset_Password{}
	type Filter struct {
		UUID string `json:"uuid"`
	}
	err := uuid_reset_password.FilterSearchSingle(&Filter{
		UUID: reset.UUID,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("UUID : %v tidak ditemukan", reset.UUID))
	}

	if uuid_reset_password.Used == true {
		return returnInvalidResponse(http.StatusNotFound, "", "Anda telah melakukan pergantian password")
	}

	diff := now.Sub(uuid_reset_password.Expired)
	if diff > 0 {
		return returnInvalidResponse(http.StatusNotFound, "", "Link Expired")
	}
	//check Borrower ID
	borrowerModel := models.Borrower{}
	err = borrowerModel.FindbyID(int(uuid_reset_password.Borrower.Int64))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun Tidak ditemukan")
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(reset.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	//get password from users entity/table
	userBorrower := models.User{}
	err = userBorrower.FindbyBorrowerID(uint64(uuid_reset_password.Borrower.Int64))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun bukan borrower personal")
	}
	userBorrower.Password = string(passwordByte)
	err = userBorrower.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Ubah Password Gagal")
	}

	uuid_reset_password.Used = true
	uuid_reset_password.Save()

	responseBody := map[string]interface{}{
		"status":  true,
		"message": "Ubah Passord berhasil",
	}

	return c.JSON(http.StatusOK, responseBody)
}
