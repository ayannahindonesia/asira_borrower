package handlers

import (
	"asira/asira"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func ClientLogin(c echo.Context) error {
	clientConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.clients", asira.App.ENV))
	if authtoken := c.Request().Header.Get("Authorization"); authtoken == clientConf["android"].(string) {
		token, err := createJwtToken("android_client", "client")
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(500, fmt.Sprintf("%s", "error"))
		}

		jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))
		expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

		return c.JSON(http.StatusOK, map[string]interface{}{
			"token":      token,
			"expires_in": expiration.Seconds(),
		})
	}

	return echo.NewHTTPError(401, fmt.Sprintf("%s", "invalid credentials"))
}
