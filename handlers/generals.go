package handlers

import (
	"encoding/base64"
	"fmt"
	"kayacredit/kc"
	"strings"
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

func createJwtToken(authtoken string, role string) (string, error) {
	var (
		claim JWTclaims
	)
	decodeByte, err := base64.StdEncoding.DecodeString(authtoken)
	client := strings.Split(string(decodeByte), ":")

	jwtConf := kc.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", kc.App.ENV))

	claim = JWTclaims{
		client[0],
		role,
		jwt.StandardClaims{
			Id:        client[0],
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
