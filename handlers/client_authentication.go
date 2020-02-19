package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//ClientLogin clientauth
func ClientLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "ClientLogin"

	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Basic "))
	if err != nil {
		NLog("warning", LogTag, fmt.Sprintf("error client authentification : %v", err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusUnauthorized, "", "Invalid Creadentials")
	}

	auth := strings.Split(string(data), ":")
	if len(auth) < 2 {
		return returnInvalidResponse(http.StatusUnauthorized, "", "Creadentials tidak ditemukan")
	}
	type Login struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}

	clientModel := models.Client{}
	err = clientModel.FilterSearchSingle(&Login{
		Key:    auth[0],
		Secret: auth[1],
	})

	if err != nil {
		NLog("error", LogTag, "client creadentials not found", c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusUnauthorized, "", "Creadentials tidak ditemukan")
	}

	token, err := createJwtToken(clientModel.Name, clientModel.Role)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("failed creating token for client %v : %v", clientModel.Name, err), c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, "", fmt.Sprint(err))
	}

	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
