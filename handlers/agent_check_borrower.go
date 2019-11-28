package handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type Response struct {
	Status          bool     `json:"status"`
	IDAgentBorrower int64    `json:"id_agent_borrower"`
	Fields          []string `json:"fields"`
}

type Filter struct {
	IdCardNumber string `json:"idcard_number" condition:"optional"`
	TaxIDnumber  string `json:"taxid_number" condition:"optional"`
	Phone        string `json:"phone" condition:"optional"`
	Email        string `json:"email" condition:"optional"`
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
	var agentBorrower models.AgentBorrower
	err = agentBorrower.FilterSearchSingle(&Filter{
		IdCardNumber: payloadFilter.IdCardNumber,
		TaxIDnumber:  payloadFilter.TaxIDnumber,
		Phone:        payloadFilter.Phone,
		Email:        payloadFilter.Email,
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
	return c.JSON(http.StatusOK, &Response{
		IDAgentBorrower: int64(agentBorrower.ID),
		Status:          true,
		Fields:          existed,
	})
}

func existingFields(agentBorrower models.AgentBorrower, payload Payload) []string {
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
		fmt.Printf("%+v\n", check)
	}
	return exists
}

func compareReflectFieldValue(is string, isReflect reflect.Value, inReflect reflect.Value) bool {
	//ambil data
	isValue := reflect.Indirect(isReflect).FieldByName(is)
	inValue := reflect.Indirect(inReflect).FieldByName(is)

	//cek equality
	// if reflect.DeepEqual(isValue, inValue) {
	if isValue.String() == inValue.String() {
		return true
	}
	return false
}
