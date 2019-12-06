package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
	"gitlab.com/asira-ayannah/basemodel"
)

func BorrowerLoanApply(c echo.Context) error {
	defer c.Request().Body.Close()
	var err error

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loan.Borrower = sql.NullInt64{Int64: int64(borrowerID), Valid: true}

	payloadRules := govalidator.MapData{
		"loan_amount":       []string{"required"},
		"installment":       []string{"required"},
		"loan_intention":    []string{"required"},
		"intention_details": []string{"required"},
		"product":           []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &loan)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = validateLoansProduct(loan)
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "validation error")
	}

	err = loan.Create()
	if err != nil {
		log.Printf("apply : %v", loan)
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Loan")
	}

	return c.JSON(http.StatusCreated, loan)
}

func BorrowerLoanGet(c echo.Context) error {
	defer c.Request().Body.Close()

	db := asira.App.DB

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	type Loans struct {
		models.Loan
		ProductName string `json:"product_name"`
		ServiceName string `json:"service_name"`
	}
	var results []Loans
	var totalRows int
	var offset int
	var rows int
	var page int

	// pagination parameters
	if c.QueryParam("rows") != "all" {
		rows, _ = strconv.Atoi(c.QueryParam("rows"))
		page, _ = strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		if rows <= 0 {
			rows = 25
		}
		offset = (page * rows) - rows
	}

	db = db.Table("loans l").
		Select("*, bp.name as product_name, bs.name as service_name").
		Joins("INNER JOIN products bp ON bp.id = l.product").
		Joins("INNER JOIN services bs ON bs.id = bp.service_id").
		Where("l.borrower = ?", borrowerID)

	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("l.status = ?", status)
	}

	if rows > 0 && offset > 0 {
		db = db.Limit(rows).Offset(offset)
	}
	err := db.Find(&results).Count(&totalRows).Error
	if err != nil {
		returnInvalidResponse(http.StatusInternalServerError, err, "pencarian loan gagal")
	}

	lastPage := int(math.Ceil(float64(totalRows) / float64(rows)))

	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        results,
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
		ID       int           `json:"id"`
		Borrower sql.NullInt64 `json:"borrower"`
	}
	err = loan.FilterSearchSingle(&Filter{
		ID: loan_id,
		Borrower: sql.NullInt64{
			Int64: int64(borrowerID),
			Valid: true,
		},
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan id %v tidak ditemukan", loan_id))
	}

	return c.JSON(http.StatusOK, loan)
}

func BorrowerLoanOTPrequest(c echo.Context) error {
	defer c.Request().Body.Close()

	loan := models.Loan{}
	borrower := models.Borrower{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loan_id, err := strconv.Atoi(c.Param("loan_id"))
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing loan id to integer")
	}

	type Filter struct {
		ID       int           `json:"id"`
		Borrower sql.NullInt64 `json:"borrower"`
	}
	err = loan.FilterSearchSingle(&Filter{
		ID: loan_id,
		Borrower: sql.NullInt64{
			Int64: int64(borrowerID),
			Valid: true,
		},
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "query result error")
	}

	borrower.FindbyID(borrowerID)

	catenate := strconv.Itoa(borrowerID) + strconv.Itoa(int(loan.ID)) // combine borrower id with loan id as counter
	counter, _ := strconv.Atoi(catenate)
	otpCode := asira.App.OTP.HOTP.At(int(counter))

	message := fmt.Sprintf("Code OTP Pengajuan Pinjaman anda adalah %s", otpCode)
	err = asira.App.Messaging.SendSMS(borrower.Phone, message)
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending otp")
	}

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
		ID       int           `json:"id"`
		Borrower sql.NullInt64 `json:"borrower"`
	}
	err = loan.FilterSearchSingle(&Filter{
		ID: loan_id,
		Borrower: sql.NullInt64{
			Int64: int64(borrowerID),
			Valid: true,
		},
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "ID tidak ditemukan")
	}

	if loan.OTPverified {
		return returnInvalidResponse(http.StatusBadRequest, "", fmt.Sprintf("loan %v sudah di verifikasi", loan_id))
	}

	catenate := strconv.Itoa(borrowerID) + strconv.Itoa(int(loan.ID)) // combine borrower id with loan id as counter
	counter, _ := strconv.Atoi(catenate)
	if asira.App.OTP.HOTP.Verify(LoanOTPverify.OTPcode, counter) {
		loan.OTPverified = true
		loan.Save()

		return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v verified", loan.ID)})
	}

	// bypass otp
	if asira.App.ENV == "development" && LoanOTPverify.OTPcode == "888999" {
		loan.OTPverified = true
		loan.Save()

		return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v verified", loan.ID)})
	}

	return returnInvalidResponse(http.StatusBadRequest, "", "OTP yang anda masukan salah")
}

func validateLoansProduct(l models.Loan) (err error) {
	var count int

	db := asira.App.DB

	err = db.Table("banks b").
		Select("p.id").
		Joins("INNER JOIN borrowers bo ON bo.bank = b.id").
		Joins("INNER JOIN services s ON s.id IN (SELECT UNNEST(b.services))").
		Joins("INNER JOIN products p ON p.service_id = s.id").
		Where("p.id = ?", l.Product).
		Where("bo.id = ?", l.Borrower.Int64).Count(&count).Error

	if count < 1 {
		err = fmt.Errorf("invalid product")
	}

	return err
}
