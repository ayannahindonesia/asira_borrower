package handlers

import (
	"asira/asira"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func ClientLogin(c echo.Context) error {
	defer c.Request().Body.Close()
	clientConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.clients", asira.App.ENV))
	if authtoken := c.Request().Header.Get("Authorization"); authtoken == clientConf["android"].(string) {
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
