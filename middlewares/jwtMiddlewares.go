package middlewares

import (
	"asira/asira"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetClientJWTmiddlewares(g *echo.Group, role string) {
	jwtConfig := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))

	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(jwtConfig["jwt_secret"].(string)),
	}))

	switch role {
	case "client":
		g.Use(validateJWTclient)
		break
	case "borrower":
		g.Use(validateJWTborrower)
		g.Use(validateBorrowerOTPstatus)
		break
	case "borrower_unverified":
		g.Use(validateJWTborrower)
		break
	}
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

		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
	}
}

func validateJWTborrower(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["role"] == "borrower" || claims["role"] == "borrower_unverified" {
				return next(c)
			} else {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid role"))
			}
		}

		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
	}
}

func validateBorrowerOTPstatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claims["otp_status"] == false {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "account needs to verify otp"))
			}

			return next(c)
		}

		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
	}
}
