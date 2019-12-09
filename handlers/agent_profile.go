package handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

//AgentProfile get current agent's profile
func AgentProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	agentModel := models.Agent{}

	agentID, _ := strconv.Atoi(claims["jti"].(string))
	err := agentModel.FindbyID(agentID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	return c.JSON(http.StatusOK, agentModel)
}

//AgentProfileEdit update current agent's profile
func AgentProfileEdit(c echo.Context) error {
	defer c.Request().Body.Close()

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	agentModel := models.Agent{}

	agentID, _ := strconv.Atoi(claims["jti"].(string))
	err := agentModel.FindbyID(agentID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	validate := validateRequestPayload(c, payloadRules, &borrowerModel)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	agentModel.Password := password

	payloadRules := govalidator.MapData{
		"name":  []string{},
		"email": []string{"unique_edit:agents,email"},
		"phone": []string{"phone", "unique:agents,taxid_number"},
	}
	return c.JSON(http.StatusOK, agentModel)
}
