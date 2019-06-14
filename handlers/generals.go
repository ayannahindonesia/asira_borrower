package handlers

import (
	"asira/asira"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type (
	JWTclaims struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		jwt.StandardClaims
	}
)

func createJwtToken(id string, role string) (string, error) {
	jwtConf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", asira.App.ENV))

	claim := JWTclaims{
		id,
		role,
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(time.Duration(jwtConf["duration"].(int)) * time.Minute).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err := rawToken.SignedString([]byte(jwtConf["jwt_secret"].(string)))
	if err != nil {
		return "", err
	}

	return token, nil
}
