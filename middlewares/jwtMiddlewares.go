package middlewares

import (
	"asira/asira"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetClientJWTmiddlewares(g *echo.Group) {
	jwtConfig := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))

	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(jwtConfig["jwt_secret"].(string)),
	}))

	g.Use(validateJWTclient)
}

func validateJWTclient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["role"] == "client" {
				return next(c)
			} else {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid role"))
			}
		}
		log.Println(token.Claims)
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
	}
}
