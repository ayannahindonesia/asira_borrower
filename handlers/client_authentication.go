package handlers

import (
	"asira/asira"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func ClientLogin(c echo.Context) error {
	defer c.Request().Body.Close()
	clientConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.clients", asira.App.ENV))
	if authtoken := strings.Trim(c.Request().Header.Get("Authorization"), "Basic "); authtoken == clientConf["admin"].(string) {
		token, err := createJwtToken("adm", "admin")
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
	if authtoken := strings.Trim(c.Request().Header.Get("Authorization"), "Basic "); authtoken == clientConf["android"].(string) {
		token, err := createJwtToken("android_client", "client")
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

	return returnInvalidResponse(http.StatusUnauthorized, "", "invalid credentials")
}
