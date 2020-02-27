package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/middlewares"
	"asira_borrower/models"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

//AgentLoanApply Agent apply new Loan
func AgentLoanApply(c echo.Context) error {
	defer c.Request().Body.Close()
	var err error

	LogTag := "AgentLoanApply"

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	loan := models.Loan{}

	payloadRules := govalidator.MapData{
		"borrower":          []string{"required"},
		"loan_amount":       []string{"required"},
		"installment":       []string{"required"},
		"loan_intention":    []string{"required"},
		"intention_details": []string{"required"},
		"product":           []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &loan)
	if validate != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "error validation",
			NLOGERR:   err,
			"payload": loan}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//is valid agent
	agent := models.Agent{}
	err = agent.FindbyID(agentID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:    "not valid agent",
			NLOGERR:    err,
			"agent_id": agentID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnauthorized, validate, "validation error : not valid agent")
	}

	//validate is borrower registered by agent
	check := agent.CheckBorrowerOwnedByAgent(loan.Borrower)
	if !check {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "not valid agent's borrower",
			NLOGERR:       err,
			"borrower_id": loan.Borrower}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnauthorized, validate, "validation error : not valid agent's borrower")
	}

	err = validateAgentLoansProduct(loan, agentID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: "not valid product",
			NLOGERR: err,
			"loan":  loan}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "validation error : not valid product")
	}

	err = loan.Create()
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error create loan",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr(),
			"loan":    loan}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Loan")
	}
	err = middlewares.SubmitKafkaPayload(loan, "loan_create")
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: "error kafka sumbit create loan",
			NLOGERR: err,
			"loan":  loan}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Sinkronisasi Borrower Baru Gagal")
	}

	NAudittrail(models.Loan{}, loan, token, "loan", fmt.Sprint(loan.ID), "agent loan apply", "agent")

	return c.JSON(http.StatusCreated, loan)
}

//AgentLoanGet Agent get his loans
func AgentLoanGet(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentLoanGet"

	db := asira.App.DB

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.Atoi(claims["jti"].(string))

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
		Joins("INNER JOIN borrowers br ON br.id = loans.borrower").
		Joins("INNER JOIN agents ag ON ag.id = br.agent_referral")

	//do join for banks
	bankID, _ := strconv.Atoi(c.QueryParam("bank"))
	if bankID > 0 {
		db = db.Joins("INNER JOIN banks bnk ON bs.id IN (SELECT UNNEST(bnk.services))  AND br.bank = bnk.id").
			Where("bnk.id = ?", bankID)
	}

	db = db.Where("ag.id = ?", agentID)

	if status := c.QueryParam("status"); len(status) > 0 {
		db = db.Where("loans.status = ?", status)
	}

	if rows > 0 && offset > 0 {
		db = db.Limit(rows).Offset(offset)
	}
	err := db.Find(&results).Count(&totalRows).Error
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:      "empty loan",
			NLOGERR:      err,
			NLOGQUERY:    asira.App.DB.QueryExpr(),
			"loan_count": totalRows}, c.Get("user").(*jwt.Token), "", false, "agent")

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

//AgentLoanGetDetails Agent get detail loan
func AgentLoanGetDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentLoanGet"

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	//cek loan
	loan := models.Loan{}
	loanID, err := strconv.ParseUint(c.Param("loan_id"), 10, 64)
	err = loan.FindbyID(loanID)
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "not valid loan",
			NLOGERR:   err,
			"loan_id": loanID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("loan id %v tidak ditemukan", loanID))
	}

	//is valid agent
	agent := models.Agent{}
	err = agent.FindbyID(agentID)
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG: "not valid agent",
			NLOGERR: err}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnauthorized, err, "validation error : not valid agent")
	}

	//validate is borrower registered by agent
	check := agent.CheckBorrowerOwnedByAgent(loan.Borrower)
	if !check {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:       "not valid agent's borrower",
			NLOGERR:       err,
			"borrower_id": loan.Borrower}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnauthorized, check, "validation error : not valid agent's borrower")
	}

	return c.JSON(http.StatusOK, loan)
}

//AgentLoanOTPrequest request for one time password
func AgentLoanOTPrequest(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentLoanOTPrequest"

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	//cek loan
	loan := models.Loan{}
	loanID, err := strconv.ParseUint(c.Param("loan_id"), 10, 64)
	err = loan.FindbyID(loanID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "not valid loan",
			NLOGERR:   err,
			"loan_id": loanID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, fmt.Sprintf("loan id %v tidak ditemukan", loanID))
	}

	//is valid agent
	agent := models.Agent{}
	err = agent.FindbyID(agentID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:    "not valid agent",
			NLOGERR:    err,
			"agent_id": agentID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnauthorized, err, "validation error : not valid agent")
	}

	//validate is borrower registered by agent
	check := agent.CheckBorrowerOwnedByAgent(loan.Borrower)
	if !check {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "not valid agent's borrower",
			NLOGERR:       err,
			"borrower_id": loan.Borrower}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnauthorized, check, "validation error : not valid agent's borrower")
	}

	//get borrower (phone)
	borrower := models.Borrower{}
	err = borrower.FindbyID(loan.Borrower)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "not valid borrower",
			NLOGERR:       err,
			"borrower_id": loan.Borrower}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "validation error : not valid borrower")
	}
	catenate := strconv.Itoa(int(loan.Borrower)) + strconv.Itoa(int(loan.ID)) // combine borrower id with loan id as counter
	counter, _ := strconv.Atoi(catenate)
	otpCode := asira.App.OTP.HOTP.At(int(counter))

	message := fmt.Sprintf("Code OTP Pengajuan Pinjaman anda adalah %s", otpCode)
	err = asira.App.Messaging.SendSMS(agent.Phone, message)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "error failed Sending SMS OTP",
			NLOGERR:   err,
			"payload": message}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending otp")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Terkirim"})
}

//AgentLoanOTPverify verify for one time password has sent
func AgentLoanOTPverify(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentLoanOTPverify"

	var LoanOTPverify struct {
		OTPcode string `json:"otp_code"`
	}

	payloadRules := govalidator.MapData{
		"otp_code": []string{"numeric", "required"},
	}

	validate := validateRequestPayload(c, payloadRules, &LoanOTPverify)
	if validate != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "error validation",
			NLOGERR:       validate,
			"otp_payload": LoanOTPverify}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)

	//cek loan ID
	loan := models.Loan{}
	loanID, err := strconv.ParseUint(c.Param("loan_id"), 10, 64)
	err = loan.FindbyID(loanID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "not valid loan",
			NLOGERR:   err,
			"loan_id": loanID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "loan id not valid")
	}

	//is valid agent
	agent := models.Agent{}
	err = agent.FindbyID(agentID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:    "not valid loan",
			NLOGERR:    err,
			"agent_id": agentID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnauthorized, err, "validation error : not valid agent")
	}

	//validate is borrower registered by agent
	check := agent.CheckBorrowerOwnedByAgent(loan.Borrower)
	if !check {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:       "not valid agent's borrower",
			NLOGERR:       "borrower not owned by agent",
			"borrower_id": loan.Borrower}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnauthorized, check, "validation error : not valid agent's borrower")
	}

	if loan.OTPverified {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: "loan already verified",
			"loan":  loan}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusBadRequest, "", fmt.Sprintf("loan %v sudah di verifikasi", loanID))
	}

	catenate := strconv.Itoa(int(loan.Borrower)) + strconv.Itoa(int(loan.ID)) // combine borrower id with loan id as counter
	counter, _ := strconv.Atoi(catenate)
	if asira.App.OTP.HOTP.Verify(LoanOTPverify.OTPcode, counter) {
		loan.OTPverified = true
		// loan.Save()
		err = middlewares.SubmitKafkaPayload(loan, "loan_update")
		if err != nil {
			NLog("error", LogTag, map[string]interface{}{
				NLOGMSG: "error kafka submit update loan",
				NLOGERR: err,
				"loan":  loan}, c.Get("user").(*jwt.Token), "", false, "agent")

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
				"loan":  loan}, c.Get("user").(*jwt.Token), "", false, "agent")

			return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal mensinkronisasi Loan")
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("loan %v verified", loan.ID)})
	}

	NLog("error", LogTag, map[string]interface{}{
		NLOGMSG:       "cannot validate OTP",
		NLOGERR:       "invalid OTP code",
		"otp_payload": LoanOTPverify}, c.Get("user").(*jwt.Token), "", false, "agent")

	return returnInvalidResponse(http.StatusBadRequest, "", "OTP yang anda masukan salah")
}

//validateLoansProduct validate product before apply loan
func validateAgentLoansProduct(l models.Loan, agentID uint64) (err error) {
	var count int

	db := asira.App.DB

	err = db.Table("banks").
		Select("p.id").
		Joins("INNER JOIN agents ag ON banks.id IN (SELECT UNNEST(ag.banks))").
		Joins("INNER JOIN services s ON s.id IN (SELECT UNNEST(banks.services))").
		Joins("INNER JOIN products p ON p.service_id = s.id").
		Where("p.id = ?", l.Product).
		Where("ag.id = ?", agentID).Count(&count).Error

	if count < 1 {
		err = fmt.Errorf("invalid product")
	}

	return err
}
