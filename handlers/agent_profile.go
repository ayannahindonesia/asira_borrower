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

	//securing old password
	password := agentModel.Password

	payloadRules := govalidator.MapData{
		"name":  []string{},
		"email": []string{"unique_edit:agents,email"},
		"phone": []string{"phone", "unique:agents,taxid_number"},
	}

	//validate request data
	validate := validateRequestPayload(c, payloadRules, &agentModel)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//restoring old password and update data
	agentModel.Password = password
	err = agentModel.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "Gagal mengubah data akun agen")
	}

	return c.JSON(http.StatusOK, agentModel)
}
