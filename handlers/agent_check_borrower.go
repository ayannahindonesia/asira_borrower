package handlers

import (
	"asira_borrower/asira"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type Response struct {
	Status          bool     `json:"status"`
	IDAgentBorrower uint64   `json:"id_agent_borrower"`
	Fields          []string `json:"fields"`
}

type Filter struct {
	IdCardNumber  string        `json:"idcard_number" condition:"optional"`
	TaxIDnumber   string        `json:"taxid_number" condition:"optional"`
	Phone         string        `json:"phone" condition:"optional"`
	Email         string        `json:"email" condition:"optional"`
	AgentReferral sql.NullInt64 `json:"agent_referral" condition:"optional"`
}

type Payload struct {
	IdCardNumber string `json:"idcard_number"`
	TaxIDnumber  string `json:"taxid_number"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

func AgentCheckBorrower(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentCheckBorrower"

	//validate post
	payloadFilter := Payload{}
	rules := govalidator.MapData{
		"idcard_number": []string{"required"},
		"taxid_number":  []string{},
		"phone":         []string{"id_phonenumber"},
		"email":         []string{"email"},
	}
	validate := validateRequestPayload(c, rules, &payloadFilter)
	if validate != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG: "error validation",
			NLOGERR: validate}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "invalid post body")
	}

	//check manual fields if not unique
	var fields = map[string]string{
		"phone":        payloadFilter.Phone,
		"email":        payloadFilter.Email,
		"taxid_number": payloadFilter.TaxIDnumber,
	}
	foundFields, err := checkUniqueFields(payloadFilter.IdCardNumber, fields)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "data already exist",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr(),,
			"fields_found": foundFields}, c.Get("user").(*jwt.Token), "", false, "agent")
		
		return returnInvalidResponse(http.StatusInternalServerError, err, "data sudah ada sebelumnya : "+foundFields)
	}

	db := asira.App.DB
	var count int

	//max borrower duplicate just == 1
	db = db.Table("borrowers").
		Select("*").
		Where("idcard_number = ? AND agent_referral <> 0", payloadFilter.IdCardNumber).
		Where(generateDeleteCheck("borrowers"))
	err = db.Count(&count).Error
	if count >= 1 {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "borrower already registered",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr(),
			"count": count}, c.Get("user").(*jwt.Token), "", false, "agent")
		
		NLog("warning", LogTag, fmt.Sprintf("borrower already registered : %v", count), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusInternalServerError, err, "borrower sudah terdaftar")
	}

	responseBody := map[string]interface{}{
		"status":  true,
		"message": "Ok",
	}
	return c.JSON(http.StatusOK, responseBody)
}

// func existingFields(agentBorrower models.Borrower, payload Payload) []string {
// 	var exists []string
// 	valPayload := reflect.ValueOf(payload)
// 	valAgentBorrower := reflect.ValueOf(agentBorrower)

// 	//loop per field agent borrower
// 	for i := 0; i < valAgentBorrower.NumField(); i++ {
// 		field := valAgentBorrower.Type().Field(i).Name
// 		//cek availability
// 		check := compareReflectFieldValue(field, valPayload, valAgentBorrower)
// 		if check == true {
// 			word, ok := EnglishToIndonesiaFields[field]
// 			if !ok {
// 				fmt.Println(err)
// 				word = field
// 			}
// 			word = strings.TrimSpace(word)
// 			fmt.Println("word =>>(", word, ")")
// 			exists = append(exists, word)
// 			fmt.Printf("%+v\n", exists)
// 		}
// 		//fmt.Printf("%+v\n", check)
// 	}
// 	return exists
// }

// func compareReflectFieldValue(is string, isReflect reflect.Value, inReflect reflect.Value) bool {
// 	//ambil data
// 	isValue := reflect.Indirect(isReflect).FieldByName(is)
// 	inValue := reflect.Indirect(inReflect).FieldByName(is)

// 	//cek equality
// 	// if reflect.DeepEqual(isValue, inValue)
// 	isVal := isValue.String()
// 	if len(isVal) > 0 && isVal == inValue.String() {
// 		return true
// 	}
// 	return false
// }
