package handlers

import (
	"fmt"
	"kayacredit/kc"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func ClientLogin(c echo.Context) error {
	clientConf := kc.App.Config.GetStringMap(fmt.Sprintf("%s.clients", kc.App.ENV))
	if authtoken := c.Request().Header.Get("Authorization"); authtoken == clientConf["android"].(string) {
		token, err := createJwtToken("android_client", "client")
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(500, fmt.Sprintf("%s", "error"))
		}

		jwtConf := kc.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", kc.App.ENV))
		expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

		return c.JSON(http.StatusOK, map[string]interface{}{
			"token":      token,
			"expires_in": expiration.Seconds(),
		})
	}

	return echo.NewHTTPError(401, fmt.Sprintf("%s", "invalid credentials"))
}
