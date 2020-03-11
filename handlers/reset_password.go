package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/crypto/bcrypt"

	guuid "github.com/google/uuid"
)

//ClientResetPassword send request reset password to email
func ClientResetPassword(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "ClientResetPassword"

	borrower := models.Borrower{}
	// r := c.Request()
	//c.Scheme() + r.Host
	baseURL := asira.App.Config.GetString(fmt.Sprintf("%s.baseurl", asira.App.ENV))
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
			NLog("error", LogTag, fmt.Sprintf("error create UUID reset password : %v", err), c.Get("user").(*jwt.Token), "", true, "")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal menyimpan UUID")
		}
		if borrower.OTPverified != true {
			NLog("error", LogTag, fmt.Sprintf("error borrower not verified : %v", borrower), c.Get("user").(*jwt.Token), "", true, "")

			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Akun Anda Belum di Verifikasi, Silahkan Verifikasi Terlebih Dahulu")
		}
		if borrower.Email == "" {
			NLog("error", LogTag, fmt.Sprintf("error not valid email for reset password : %v", borrower), c.Get("user").(*jwt.Token), "", true, "")

			return returnInvalidResponse(http.StatusUnprocessableEntity, "", "Email Tidak Ditemukan atau Email anda belum terdaftar")
		}

		//send email
		to := borrower.Email
		subject := "[NO REPLY] - Reset Password Aplikasi ASIRA"
		link := baseURL + "/deepLinks/" + uuid.UUID
		message := "Hai Nasabah,\n\nIni adalah email untuk melakukan reset login akun anda. Silahkan klik link di bawah ini agar dapat melakukan reset login akun.\nLink ini hanya valid dalam waktu 24 jam.\n" + link + " \n\n\n Ayannah Solusi Nusantara Team"

		err = asira.App.Emailer.SendMail(to, subject, message)
		if err != nil {
			NLog("error", LogTag, fmt.Sprintf("error send Email reset password : %v", err), c.Get("user").(*jwt.Token), "", true, "")

			log.Println(err.Error())
		}

		NLog("event", LogTag, fmt.Sprintf("success reset password : %v", borrower.Email), c.Get("user").(*jwt.Token), "", true, "")

		return c.JSON(http.StatusOK, map[string]interface{}{"message": "Link reset password telah dikirim ke email, silahkan cek email anda"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Link reset password telah dikirim ke email, silahkan cek email anda"})
}

//ChangePassword change password with UUID get from ClientResetPassword(...)
func ChangePassword(c echo.Context) error {
	defer c.Request().Body.Close()
	now := time.Now()

	LogTag := "ChangePassword"

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
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: "error validation",
			NLOGERR: validate}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//check UUID
	clientResetPassword := models.Uuid_Reset_Password{}
	type Filter struct {
		UUID string `json:"uuid"`
	}
	err := clientResetPassword.FilterSearchSingle(&Filter{
		UUID: reset.UUID,
	})
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: fmt.Sprintf("error invalid UUID : %v", reset.UUID),
			NLOGERR: err}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("UUID : %v tidak ditemukan", reset.UUID))
	}

	if clientResetPassword.Used == true {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: fmt.Sprintf("already used UUID : %v", reset.UUID),
			NLOGERR: "clientResetPassword.Used == true"}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusNotFound, "", "Anda telah melakukan pergantian password")
	}

	diff := now.Sub(clientResetPassword.Expired)
	if diff > 0 {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: fmt.Sprintf("expired used UUID : %v", reset.UUID),
			NLOGERR: diff}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusNotFound, "", "Link Expired")
	}
	//check Borrower ID
	borrowerModel := models.Borrower{}
	err = borrowerModel.FindbyID(uint64(clientResetPassword.Borrower.Int64))
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "borrower not found",
			NLOGERR:       err,
			"borrower_id": clientResetPassword.Borrower.Int64}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun Tidak ditemukan")
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(reset.Password), bcrypt.DefaultCost)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "error generate password hash",
			NLOGERR:       err,
			"borrower_id": clientResetPassword.Borrower.Int64}, c.Get("user").(*jwt.Token), "", true, "")

		return err
	}

	//get password from users entity/table
	userBorrower := models.User{}
	err = userBorrower.FindbyBorrowerID(uint64(clientResetPassword.Borrower.Int64))
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "not valid personal borrower",
			NLOGERR:       err,
			"borrower_id": clientResetPassword.Borrower.Int64}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun bukan borrower personal")
	}
	userBorrower.Password = string(passwordByte)
	err = userBorrower.Save()
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "error update user borrower password ",
			NLOGERR:       err,
			"borrower_id": clientResetPassword.Borrower.Int64}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Ubah Password Gagal")
	}

	clientResetPassword.Used = true
	clientResetPassword.Save()

	responseBody := map[string]interface{}{
		"status":  true,
		"message": "Ubah Passord berhasil",
	}

	NLog("info", LogTag, map[string]interface{}{
		NLOGMSG:       "success change password",
		"borrower_id": clientResetPassword.Borrower.Int64}, c.Get("user").(*jwt.Token), "", true, "")

	return c.JSON(http.StatusOK, responseBody)
}
