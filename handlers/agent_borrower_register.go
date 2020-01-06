package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type (
	FilterAgentRef struct {
		ID            uint64 `json:"id"`
		AgentReferral int64  `json:"agent_referral"`
	}

	FilterAgentPhone struct {
		ID    uint64 `json:"id"`
		Phone string `json:"phone"`
	}
)

//AgentRegisterBorrower agent register new borrower
func AgentRegisterBorrower(c echo.Context) error {
	defer c.Request().Body.Close()
	type (
		Register struct {
			Fullname             string    `json:"fullname"`
			Nickname             string    `json:"nickname"`
			Gender               string    `json:"gender" `
			Image                string    `json:"image"`
			IdCardNumber         string    `json:"idcard_number" `
			IdCardImage          string    `json:"idcard_image"`
			TaxIDImage           string    `json:"taxid_image"`
			TaxIDnumber          string    `json:"taxid_number"`
			Nationality          string    `json:"nationality"`
			Email                string    `json:"email"`
			Birthday             time.Time `json:"birthday"`
			Birthplace           string    `json:"birthplace"`
			LastEducation        string    `json:"last_education"`
			MotherName           string    `json:"mother_name"`
			Phone                string    `json:"phone"`
			MarriedStatus        string    `json:"marriage_status"`
			SpouseName           string    `json:"spouse_name"`
			SpouseBirthday       time.Time `json:"spouse_birthday"`
			SpouseLastEducation  string    `json:"spouse_lasteducation"`
			Dependants           int       `json:"dependants,omitempty"`
			Address              string    `json:"address"`
			Province             string    `json:"province"`
			City                 string    `json:"city"`
			NeighbourAssociation string    `json:"neighbour_association"`
			Hamlets              string    `json:"hamlets"`
			HomePhoneNumber      string    `json:"home_phonenumber"`
			Subdistrict          string    `json:"subdistrict"`
			UrbanVillage         string    `json:"urban_village"`
			HomeOwnership        string    `json:"home_ownership"`
			LivedFor             int       `json:"lived_for"`
			Occupation           string    `json:"occupation"`
			EmployeeID           string    `json:"employee_id"`
			EmployerName         string    `json:"employer_name"`
			EmployerAddress      string    `json:"employer_address"`
			Department           string    `json:"department"`
			BeenWorkingFor       int       `json:"been_workingfor"`
			DirectSuperior       string    `json:"direct_superiorname"`
			EmployerNumber       string    `json:"employer_number"`
			MonthlyIncome        int       `json:"monthly_income"`
			OtherIncome          int       `json:"other_income"`
			OtherIncomeSource    string    `json:"other_incomesource"`
			FieldOfWork          string    `json:"field_of_work"`
			RelatedPersonName    string    `json:"related_personname"`
			RelatedRelation      string    `json:"related_relation"`
			RelatedPhoneNumber   string    `json:"related_phonenumber"`
			RelatedHomePhone     string    `json:"related_homenumber"`
			RelatedAddress       string    `json:"related_address"`
			Bank                 int64     `json:"bank"`
			BankAccountNumber    string    `json:"bank_accountnumber"`
		}
	)
	register := Register{}

	payloadRules := govalidator.MapData{
		"fullname":              []string{"required"},
		"nickname":              []string{},
		"gender":                []string{"required"},
		"image":                 []string{},
		"idcard_number":         []string{"required"},
		"taxid_number":          []string{},
		"nationality":           []string{},
		"email":                 []string{"email"},
		"birthday":              []string{"date"},
		"birthplace":            []string{"required"},
		"last_education":        []string{"required"},
		"mother_name":           []string{"required"},
		"phone":                 []string{"id_phonenumber"},
		"marriage_status":       []string{"required"},
		"spouse_name":           []string{},
		"spouse_birthday":       []string{"date"},
		"spouse_lasteducation":  []string{},
		"dependants":            []string{},
		"address":               []string{"required"},
		"province":              []string{"required"},
		"city":                  []string{"required"},
		"neighbour_association": []string{"required"},
		"hamlets":               []string{"required"},
		"home_phonenumber":      []string{},
		"subdistrict":           []string{"required"},
		"urban_village":         []string{"required"},
		"home_ownership":        []string{"required"},
		"lived_for":             []string{"required"},
		"occupation":            []string{"required"},
		"employee_id":           []string{},
		"employer_name":         []string{"required"},
		"employer_address":      []string{"required"},
		"department":            []string{"required"},
		"been_workingfor":       []string{},
		"direct_superiorname":   []string{},
		"employer_number":       []string{"required"},
		"monthly_income":        []string{"required"},
		"other_income":          []string{},
		"other_incomesource":    []string{},
		"field_of_work":         []string{"required"},
		"related_personname":    []string{"required"},
		"related_relation":      []string{"required"},
		"related_phonenumber":   []string{"required"},
		"related_homenumber":    []string{},
		"bank":                  []string{"required"},
		"bank_accountnumber":    []string{},
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentModel := models.Agent{}
	agentID, _ := strconv.ParseInt(claims["jti"].(string), 10, 64)
	err := agentModel.FindbyID(int(agentID))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun Agen tidak ditemukan")
	}

	//validate
	validate := validateRequestPayload(c, payloadRules, &register)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek bank di dalam list bank agent atau tidak
	validBank := false
	for _, val := range agentModel.Banks {
		if register.Bank == val {
			validBank = true
			break
		}
	}
	if !validBank {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Bank tidak terdaftar untuk agent")
	}

	r, err := json.Marshal(register)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}
	borrower := models.Borrower{}
	json.Unmarshal(r, &borrower)

	//upload image profile borrower
	ImageProfil := ""
	if register.Image != "" || len(register.Image) != 0 {
		ImageProfil, err = uploadImageS3Formatted("boragn", register.Image)
		if err != nil {
			return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal : Image profil failed to upload")
		}
	}

	//upload image id card
	IdCardImage, err := uploadImageS3Formatted("ktp", register.IdCardImage)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal : IDCardImage failed to upload")
	}

	//upload image tax card
	TaxIDImage, err := uploadImageS3Formatted("tax", register.TaxIDImage)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal : TaxIDImage failed to upload")
	}

	//encrypted
	encryptPassphrase := asira.App.Config.GetString(fmt.Sprintf("%s.passphrase", asira.App.ENV))
	borrower.IdCardImage, err = encrypt(IdCardImage, encryptPassphrase)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi Id card gagal")
	}
	borrower.TaxIDImage, err = encrypt(TaxIDImage, encryptPassphrase)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Enkripsi NPWP gagal")
	}

	//set vars
	if ImageProfil != "" {
		borrower.Image = ImageProfil
	}

	borrower.AgentReferral = sql.NullInt64{
		Int64: agentID,
		Valid: true,
	}
	borrower.Bank = sql.NullInt64{
		Int64: int64(register.Bank),
		Valid: true,
	}

	//check manual fields if not unique
	var fields = map[string]string{
		"phone":              register.Phone,
		"email":              register.Email,
		"taxid_number":       register.TaxIDnumber,
		"bank_accountnumber": register.BankAccountNumber,
	}
	foundFields, err := checkUniqueFields(register.IdCardNumber, fields)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "data sudah ada sebelumnya : "+foundFields)
	}

	//max borrower duplicate just == 1
	db := asira.App.DB
	var count int
	db = db.Table("borrowers").
		Select("*").
		Where("idcard_number = ? AND agent_referral <> 0", register.IdCardNumber)
	err = db.Count(&count).Error
	if count >= 1 {
		return returnInvalidResponse(http.StatusInternalServerError, err, "borrower sudah terdaftar")
	}

	r, err := json.Marshal(register)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}
	json.Unmarshal(r, &borrower)

	//set need to OTP verify and create new borrower
	borrower.OTPverified = false
	err = borrower.Create()
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "Pendaftaran Borrower Baru Gagal")
	}

	return c.JSON(http.StatusCreated, borrower)
}

//AgentRequestOTP request OTP for after registered new borrower
func AgentRequestOTP(c echo.Context) error {
	defer c.Request().Body.Close()

	otpRequest := VerifyAccountOTPrequest{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	AgentID, _ := strconv.ParseInt(claims["jti"].(string), 10, 64)
	borrowerID, _ := strconv.ParseUint(c.Param("borrower_id"), 10, 64)

	//cek borrower owned by agent
	borrower := models.Borrower{}
	err := borrower.FilterSearchSingle(&FilterAgentRef{
		ID:            borrowerID,
		AgentReferral: AgentID,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, err, "validation error : not valid agent's borrower")
	}

	payloadRules := govalidator.MapData{
		"phone": []string{"regex:^[0-9]+$", "required"},
	}

	validate := validateRequestPayload(c, payloadRules, &otpRequest)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek agent's phone
	agent := models.Agent{}
	err = agent.FilterSearchSingle(&FilterAgentPhone{
		ID:    uint64(AgentID),
		Phone: otpRequest.Phone,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, err, "validation error : not valid agent's phone")
	}

	//generate OTP
	catenate := strconv.Itoa(int(borrowerID)) + agent.Phone[len(agent.Phone)-4:] // combine borrower id with last 4 digit of phone as counter
	counter, _ := strconv.Atoi(catenate)
	otpCode := asira.App.OTP.HOTP.At(int(counter))

	//Send OTP sms
	message := fmt.Sprintf("Code OTP Registrasi anda adalah %s", otpCode)
	err = asira.App.Messaging.SendSMS(agent.Phone, message)
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending otp")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Terkirim"})
}

//AgentVerifyOTP verify OTP after call request OTP (AgentRequestOTP)
func AgentVerifyOTP(c echo.Context) error {
	defer c.Request().Body.Close()

	otpVerify := VerifyAccountOTPverify{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	AgentID, _ := strconv.ParseInt(claims["jti"].(string), 10, 64)
	borrowerID, _ := strconv.ParseUint(c.Param("borrower_id"), 10, 64)

	//cek borrower owned by agent
	borrower := models.Borrower{}
	err := borrower.FilterSearchSingle(&FilterAgentRef{
		ID:            borrowerID,
		AgentReferral: AgentID,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, err, "validation error : not valid agent's borrower")
	}

	payloadRules := govalidator.MapData{
		"phone":    []string{"regex:^[0-9]+$", "required"},
		"otp_code": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &otpVerify)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//cek agent's phone
	agent := models.Agent{}
	err = agent.FilterSearchSingle(&FilterAgentPhone{
		ID:    uint64(AgentID),
		Phone: otpVerify.Phone,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, err, "validation error : not valid agent's phone")
	}

	catenate := strconv.Itoa(int(borrowerID)) + agent.Phone[len(agent.Phone)-4:] // combine borrower id with last 4 digit of phone as counter
	counter, _ := strconv.Atoi(catenate)
	if asira.App.OTP.HOTP.Verify(otpVerify.OTPcode, counter) {
		updateAgentBorrowerOTPStatus(int(borrowerID))
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Verified"})
	}

	// bypass otp
	if asira.App.ENV == "development" && otpVerify.OTPcode == "888999" {
		updateAgentBorrowerOTPStatus(int(borrowerID))
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "OTP Verified"})
	}

	return returnInvalidResponse(http.StatusBadRequest, "", "OTP salah")
}

func updateAgentBorrowerOTPStatus(borrowerID int) {
	modelBorrower := models.Borrower{}
	_ = modelBorrower.FindbyID(borrowerID)
	modelBorrower.OTPverified = true
	modelBorrower.Save()
}
