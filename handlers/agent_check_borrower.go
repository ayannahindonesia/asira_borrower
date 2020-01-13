package handlers

import (
	"asira_borrower/models"
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

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
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "invalid post body")
	}

	//check is agent's borrower exist or not
	var agentBorrower models.Borrower
	err = agentBorrower.FilterSearchSingle(&Filter{
		IdCardNumber: payloadFilter.IdCardNumber,
		TaxIDnumber:  payloadFilter.TaxIDnumber,
		Phone:        payloadFilter.Phone,
		Email:        payloadFilter.Email,
		AgentReferral: sql.NullInt64{
			Int64: 0,
			Valid: true,
		},
	})

	//if not exist yet
	if err != nil {
		return c.JSON(http.StatusOK, &Response{
			IDAgentBorrower: 0,
			Status:          false,
			Fields:          nil,
		})
	}
	//if exist
	existed := existingFields(agentBorrower, payloadFilter)

	//set
	id := agentBorrower.ID
	status := true

	//if fields duplicate not found for agent's borrower (AgentReferral != 0)
	if len(existed) == 0 {
		id = 0
		status = false

		//else if..error existed but AgentReferral == 0 (personal)
	} else if agentBorrower.AgentReferral.Int64 == 0 {

		//if found duplicate (existed) but just 1, that is IdCardNumber then skip
		if len(existed) == 1 && existed[0] == "IdCardNumber" {
			id = 0
			status = false
			existed = nil
		}
	}
	return c.JSON(http.StatusOK, &Response{
		IDAgentBorrower: id,
		Status:          status,
		Fields:          existed,
	})
}

func existingFields(agentBorrower models.Borrower, payload Payload) []string {
	var exists []string
	valPayload := reflect.ValueOf(payload)
	valAgentBorrower := reflect.ValueOf(agentBorrower)

	//loop per field agent borrower
	for i := 0; i < valAgentBorrower.NumField(); i++ {
		field := valAgentBorrower.Type().Field(i).Name
		//cek availability
		check := compareReflectFieldValue(field, valPayload, valAgentBorrower)
		if check == true {
			exists = append(exists, field)
			fmt.Printf("%+v\n", exists)
		}
		//fmt.Printf("%+v\n", check)
	}
	return exists
}

func compareReflectFieldValue(is string, isReflect reflect.Value, inReflect reflect.Value) bool {
	//ambil data
	isValue := reflect.Indirect(isReflect).FieldByName(is)
	inValue := reflect.Indirect(inReflect).FieldByName(is)

	//cek equality
	// if reflect.DeepEqual(isValue, inValue)
	isVal := isValue.String()
	if len(isVal) > 0 && isVal == inValue.String() {
		return true
	}
	return false
}
