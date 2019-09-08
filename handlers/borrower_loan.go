package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func BorrowerLoanApply(c echo.Context) error {
	defer c.Request().Body.Close()

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loan.Owner = sql.NullInt64{Int64: int64(borrowerID), Valid: true}

	payloadRules := govalidator.MapData{
		"loan_amount":       []string{"required"},
		"installment":       []string{"required"},
		"loan_intention":    []string{"required", "loan_purposes"},
		"intention_details": []string{"required"},
		"service":           []string{"required"},
		"product":           []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &loan)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	newLoan, err := loan.Create()
	if err != nil {
		log.Printf("apply : %v", loan)
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Loan")
	}

	return c.JSON(http.StatusCreated, newLoan)
}

func BorrowerLoanGet(c echo.Context) error {
	defer c.Request().Body.Close()

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	status := c.QueryParam("status")

	type Filter struct {
		Owner  sql.NullInt64 `json:"owner"`
		Status string        `json:"status"`
	}

	result, err := loan.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		Owner: sql.NullInt64{
			Int64: int64(borrowerID),
			Valid: true,
		},
		Status: status,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}

func BorrowerLoanGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loan_id, err := strconv.Atoi(c.Param("loan_id"))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Loan Id tidak ditemukan")
	}

	type Filter struct {
		ID    int           `json:"id"`
		Owner sql.NullInt64 `json:"owner"`
	}
	result, err := loan.FilterSearchSingle(&Filter{
		ID: loan_id,
		Owner: sql.NullInt64{
			Int64: int64(borrowerID),
			Valid: true,
		},
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan id %v tidak ditemukan", loan_id))
	}

	return c.JSON(http.StatusOK, result)
}

func BorrowerLoanOTPrequest(c echo.Context) error {
	defer c.Request().Body.Close()

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loan_id, err := strconv.Atoi(c.Param("loan_id"))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing loan id to integer")
	}

	type Filter struct {
		ID    int           `json:"id"`
		Owner sql.NullInt64 `json:"owner"`
	}
	result, err := loan.FilterSearchSingle(&Filter{
		ID: loan_id,
		Owner: sql.NullInt64{
			Int64: int64(borrowerID),
			Valid: true,
		},
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "query result error")
	}

	catenate := strconv.Itoa(borrowerID) + strconv.Itoa(int(result.ID)) // combine borrower id with loan id as counter
	counter, _ := strconv.Atoi(catenate)
	otpCode := asira.App.OTP.HOTP.At(int(counter))
	log.Printf("generate otp code for loan id %v : %v", result.ID, otpCode)

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Terkirim"})
}

func BorrowerLoanOTPverify(c echo.Context) error {
	defer c.Request().Body.Close()

	var LoanOTPverify struct {
		OTPcode string `json:"otp_code"`
	}

	payloadRules := govalidator.MapData{
		"otp_code": []string{"numeric", "required"},
	}

	validate := validateRequestPayload(c, payloadRules, &LoanOTPverify)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loan_id, err := strconv.Atoi(c.Param("loan_id"))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing loan id to integer")
	}

	type Filter struct {
		ID    int           `json:"id"`
		Owner sql.NullInt64 `json:"owner"`
	}
	result, err := loan.FilterSearchSingle(&Filter{
		ID: loan_id,
		Owner: sql.NullInt64{
			Int64: int64(borrowerID),
			Valid: true,
		},
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "ID tidak ditemukan")
	}

	if result.OTPverified {
		return returnInvalidResponse(http.StatusBadRequest, "", fmt.Sprintf("loan %v sudah di verifikasi", loan_id))
	}

	catenate := strconv.Itoa(borrowerID) + strconv.Itoa(int(result.ID)) // combine borrower id with loan id as counter
	counter, _ := strconv.Atoi(catenate)
	if asira.App.OTP.HOTP.Verify(LoanOTPverify.OTPcode, counter) {
		result.OTPverified = true
		result.Save()

		return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v verified", result.ID)})
	}

	// bypass otp
	if asira.App.ENV == "development" && LoanOTPverify.OTPcode == "888999" {
		result.OTPverified = true
		result.Save()

		return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v verified", result.ID)})
	}

	return returnInvalidResponse(http.StatusBadRequest, "", "OTP yang anda masukan salah")
}
