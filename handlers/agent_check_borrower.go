package handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type Response struct {
	Status          bool     `json:"status"`
	IDAgentBorrower int64    `json:"id_agent_borrower"`
	Fields          []string `json:"fields"`
}

type Filter struct {
	IDCardNumber string `json:"idcard_number" condition:"LIKE,optional"`
	TaxIDNumber  string `json:"taxid_number" condition:"LIKE,optional"`
	Phone        string `json:"phone" condition:"LIKE,optional"`
	Email        string `json:"email" condition:"LIKE,optional"`
}

type Payload struct {
	IDCardNumber string `json:"idcard_number"`
	TaxIDNumber  string `json:"taxid_number"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

func AgentCheckBorrower(c echo.Context) error {
	defer c.Request().Body.Close()

	//validate agent
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, err := strconv.ParseInt(claims["jti"].(string), 10, 64)
	var agent models.Agent
	err = agent.FindbyID(int(agentID))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

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
	err = agentBorrower.FilterSearchSingleWhereOr(&Filter{
		IDCardNumber: payloadFilter.IDCardNumber,
		TaxIDNumber:  payloadFilter.TaxIDNumber,
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
		//fmt.Printf("%+v\n", valAgentBorrower.Field(i).Interface())
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

	// if reflect.Zero(isValue).Int() == 0 || reflect.Zero(inValue).Int() == 0 {
	// 	return false
	// }
	fmt.Printf("reflect.DeepEqual(%+v, %+v) == %+v\n", isValue, inValue, reflect.DeepEqual(isValue, inValue))
	//cek equality
	// if reflect.DeepEqual(isValue, inValue) {
	if isValue.String() == inValue.String() {
		return true
	}
	return false
}
