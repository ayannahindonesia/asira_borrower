package handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

func FCMTokenUpdate(c echo.Context) error {
	defer c.Request().Body.Close()

	var payload models.User

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userBorrower := models.User{}

	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := userBorrower.FindbyBorrowerID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

	payloadRules := govalidator.MapData{
		"fcm_token": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &payload)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//update FCMToken
	userBorrower.FCMToken = payload.FCMToken
	err = userBorrower.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error saving Password")
	}

	responseBody := map[string]interface{}{
		"status":  true,
		"message": "FCM Token Updated",
	}
	return c.JSON(http.StatusOK, responseBody)
}
