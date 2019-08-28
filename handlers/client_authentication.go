package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func ClientLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Basic "))
	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Invalid Creadentials")
	}

	auth := strings.Split(string(data), ":")
	type Login struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}

	clientModel := models.Client_config{}
	clientData, err := clientModel.FilterSearchSingle(&Login{
		Key:    auth[0],
		Secret: auth[1],
	})

	if err != nil {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Credentials tidak ditemukan")
	}

	token, err := createJwtToken(clientData.Name, clientData.Role)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, "", fmt.Sprint(err))
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
