package handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

type Response struct {
	IDAgentBorrower int64 `json:"id_agent_borrower"`
	Status          bool  `json:"status"`
}

func AgentCheckBorrower(c echo.Context) error {
	defer c.Request().Body.Close()

	type Filter struct {
		AgentID      int64  `json:"agent_id" `
		IDCardNumber string `json:"idcard_number" condition:"LIKE"`
		TaxIDNumber  string `json:"taxid_number" condition:"LIKE"`
		Phone        string `json:"phone" condition:"LIKE"`
		Email        string `json:"email" condition:"LIKE"`
	}

	type Payload struct {
		IDCardNumber string `json:"idcard_number"`
		TaxIDNumber  string `json:"taxid_number"`
		Phone        string `json:"phone"`
		Email        string `json:"email"`
	}

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
		AgentID:      agentID,
		IDCardNumber: payloadFilter.IDCardNumber,
		TaxIDNumber:  payloadFilter.TaxIDNumber,
		Phone:        payloadFilter.Phone,
		Email:        payloadFilter.Email,
	})
	//if exist
	if err != nil {
		return c.JSON(http.StatusOK, &Response{
			IDAgentBorrower: 0,
			Status:          false,
		})
	}
	//if not exist yet
	return c.JSON(http.StatusOK, &Response{
		IDAgentBorrower: int64(agentBorrower.ID),
		Status:          true,
	})
}
