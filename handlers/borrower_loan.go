package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/middlewares"
	"asira_borrower/models"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

//BorrowerLoanApply borrower apply new Loan
func BorrowerLoanApply(c echo.Context) error {
	defer c.Request().Body.Close()
	var err error
	LogTag := "BorrowerLoanApply"

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	loan.Borrower = borrowerID

	payloadRules := govalidator.MapData{
		"loan_amount":       []string{"required"},
		"installment":       []string{"required"},
		"loan_intention":    []string{"required"},
		"intention_details": []string{"required"},
		"product":           []string{"required"},
		"form_info":         []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &loan)
	if validate != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "error validation ",
			NLOGERR: validate,
		}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek active loan exist or not
	if isBorrowerHaveActiveLoan(borrowerID) == true {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "error validation borrower had processing / active loan",
			NLOGERR: validate,
		}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error : sudah memiliki loan aktif atau sedang diproses sebelumnya")
	}

	err = validateLoansProduct(loan)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error validate products",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr(),
			"loan":    loan}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "validation error")
	}

	FormInfoUploadImages(&loan)

	//must create before kafka sync
	err = loan.Create()
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error create loan",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr(),
			"loan":    loan}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Loan")
	}

	err = middlewares.SubmitKafkaPayload(loan, "loan_create")
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: "error submit kafka create loan",
			NLOGERR: err,
			"loan":  loan}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Sinkronisasi Borrower Baru Gagal")
	}

	NLog("info", LogTag, map[string]interface{}{NLOGMSG: "success create loan", "loan": loan}, c.Get("user").(*jwt.Token), "", false, "borrower")

	NAudittrail(models.Loan{}, loan, token, "borrower", fmt.Sprint(loan.ID), "borrower loan apply", "borrower")

	return c.JSON(http.StatusCreated, loan)
}

//BorrowerLoanGet borrower get his loans
func BorrowerLoanGet(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "BorrowerLoanGet"

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

	db = db.Table("loans").
		Select("*, bp.name as product_name, bs.name as service_name").
		Joins("INNER JOIN products bp ON bp.id = loans.product").
		Joins("INNER JOIN services bs ON bs.id = bp.service_id").
		Where("loans.borrower = ?", borrowerID)

	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("loans.status = ?", status)
	}

	if rows > 0 && offset > 0 {
		db = db.Limit(rows).Offset(offset)
	}
	err := db.Find(&results).Count(&totalRows).Error
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "empty loan list",
			NLOGERR: err}, c.Get("user").(*jwt.Token), "", false, "borrower")

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

//BorrowerLoanGetDetails borrower get detail loan
func BorrowerLoanGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "BorrowerLoanGetDetails"

	type LoanSelect struct {
		models.Loan
		Installments []models.Installment `json:"installment_details"`
	}
	loan := LoanSelect{}
	installments := []models.Installment{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loanID, err := strconv.Atoi(c.Param("loan_id"))
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "invalid loan ID",
			NLOGERR:   err,
			"loan_id": c.Param("loan_id")}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Loan Id tidak valid")
	}

	type Filter struct {
		ID       int    `json:"id"`
		Borrower uint64 `json:"borrower"`
	}

	err = loan.FilterSearchSingle(&Filter{
		ID:       loanID,
		Borrower: uint64(borrowerID),
	})
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "not found loan",
			NLOGERR:   err,
			"loan_id": loanID}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan id %v tidak ditemukan", loanID))
	}

	db := asira.App.DB
	err = db.Table("installments").
		Select("*").
		Where("id IN (?)", strings.Fields(strings.Trim(fmt.Sprint(loan.InstallmentID), "[]"))).
		Scan(&installments).Error
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{NLOGMSG: "query not found : '%v' error : %v", "query": db.QueryExpr(), NLOGERR: err}, c.Get("user").(*jwt.Token), "", false, "agent")
	}
	loan.Installments = installments

	return c.JSON(http.StatusOK, loan)
}

//BorrowerLoanOTPrequest request for one time password
func BorrowerLoanOTPrequest(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "BorrowerLoanOTPrequest"

	loan := models.Loan{}
	borrower := models.Borrower{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	loanID, err := strconv.ParseUint(c.Param("loan_id"), 10, 64)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error parsing loan ID",
			NLOGERR:   err,
			"loan_id": c.Param("loan_id")}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing loan id to integer")
	}

	type Filter struct {
		ID       uint64 `json:"id"`
		Borrower uint64 `json:"borrower"`
	}
	err = loan.FilterSearchSingle(&Filter{
		ID:       loanID,
		Borrower: borrowerID,
	})
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "invalid loan for borrower",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr(),
			"loan_id": loanID}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusNotFound, err, "query result error")
	}

	borrower.FindbyID(borrowerID)

	catenate := strconv.Itoa(int(borrowerID)) + strconv.Itoa(int(loan.ID)) // combine borrower id with loan id as counter
	counter, _ := strconv.Atoi(catenate)
	otpCode := asira.App.OTP.HOTP.At(int(counter))

	message := fmt.Sprintf("Code OTP Pengajuan Pinjaman anda adalah %s", otpCode)
	err = asira.App.Messaging.SendSMS(borrower.Phone, message)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error failed Sending SMS OTP",
			NLOGERR:   err,
			"payload": message}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending otp")
	}

	NLog("info", LogTag, map[string]interface{}{
		NLOGMSG:   "success send OTP request verify loan",
		"loan_id": loanID,
		"phone":   borrower.Phone}, c.Get("user").(*jwt.Token), "", false, "borrower")

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Terkirim"})
}

//BorrowerLoanOTPverify verify for one time password has sent
func BorrowerLoanOTPverify(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "BorrowerLoanOTPverify"

	var LoanOTPverify struct {
		OTPcode string `json:"otp_code"`
	}

	payloadRules := govalidator.MapData{
		"otp_code": []string{"numeric", "required"},
	}

	validate := validateRequestPayload(c, payloadRules, &LoanOTPverify)
	if validate != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "invalid loan for borrower",
			NLOGERR:   validate,
			"payload": LoanOTPverify}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	loanID, err := strconv.ParseUint(c.Param("loan_id"), 10, 64)
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "error parsing loan ID",
			NLOGERR:   err,
			"loan_id": c.Param("loan_id")}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error parsing loan id to integer")
	}

	type Filter struct {
		ID       uint64 `json:"id"`
		Borrower uint64 `json:"borrower"`
	}
	err = loan.FilterSearchSingle(&Filter{
		ID:       loanID,
		Borrower: borrowerID,
	})
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "invalid loan for borrower",
			NLOGERR:   err,
			"loan_id": loanID}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusNotFound, err, "ID tidak ditemukan")
	}

	//is already verified
	if loan.OTPverified {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "loan already verified",
			"loan_id": loanID}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusBadRequest, "", fmt.Sprintf("loan %v sudah di verifikasi", loanID))
	}

	catenate := strconv.Itoa(int(borrowerID)) + strconv.Itoa(int(loan.ID)) // combine borrower id with loan id as counter
	counter, _ := strconv.Atoi(catenate)
	if asira.App.OTP.HOTP.Verify(LoanOTPverify.OTPcode, counter) {
		loan.OTPverified = true
		// loan.Save()
		err = middlewares.SubmitKafkaPayload(loan, "loan_update")
		if err != nil {
			NLog("error", LogTag, map[string]interface{}{
				NLOGMSG: "error kafka submit update loan",
				NLOGERR: err,
				"loan":  loan}, c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal mensinkronisasi Loan")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v verified", loan.ID)})
	}

	// bypass otp
	if asira.App.ENV == "development" && LoanOTPverify.OTPcode == "888999" {
		loan.OTPverified = true
		// loan.Save()
		err = middlewares.SubmitKafkaPayload(loan, "loan_update")
		if err != nil {
			NLog("error", LogTag, map[string]interface{}{
				NLOGMSG: "error kafka submit update loan",
				NLOGERR: err,
				"loan":  loan}, c.Get("user").(*jwt.Token), "", false, "borrower")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal mensinkronisasi Loan")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v verified", loan.ID)})
	}

	return returnInvalidResponse(http.StatusBadRequest, "", "OTP yang anda masukan salah")
}

//validateLoansProduct validate product before apply loan
func validateLoansProduct(l models.Loan) (err error) {
	var count int

	db := asira.App.DB

	err = db.Table("banks").
		Select("p.id").
		Joins("INNER JOIN borrowers bo ON bo.bank = banks.id").
		Joins("INNER JOIN services s ON s.id IN (SELECT UNNEST(banks.services))").
		Joins("INNER JOIN products p ON p.service_id = s.id").
		Where("p.id = ?", l.Product).
		Where("bo.id = ?", l.Borrower).Count(&count).Error

	if count < 1 {
		err = fmt.Errorf("invalid product")
	}

	return err
}

// FormInfoUploadImages func
func FormInfoUploadImages(l *models.Loan) {
	type Forminfo struct {
		Label   string `json:"label"`
		Status  string `json:"status"`
		Type    string `json:"type"`
		Value   string `json:"value"`
		Answers string `json:"answers"`
	}

	var forminfos []Forminfo
	detectimage := false

	json.Unmarshal(l.FormInfo.RawMessage, &forminfos)

	for k, v := range forminfos {
		if v.Type == "image" {
			forminfoimageurl, err := uploadImageS3Formatted(v.Label, v.Answers)
			if err != nil {
				NLog("error", "FormInfoUploadImages", map[string]interface{}{
					NLOGMSG: "error uploading Form Info image",
					NLOGERR: err,
					"loan":  l.ID}, nil, "", true, "loan")
			}

			forminfos[k].Answers = forminfoimageurl
			detectimage = true
		}
	}

	if detectimage {
		b, _ := json.Marshal(forminfos)
		l.FormInfo = postgres.Jsonb{b}
	}
}
