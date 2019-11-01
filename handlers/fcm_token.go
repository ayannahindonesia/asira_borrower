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

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrower := models.Borrower{}

	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrower.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

	payloadRules := govalidator.MapData{
		"fcm_token": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	borrower.FCMToken = borrower.FCMToken
	err = borrower.Save()
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error saving Password")
	}
	responseBody := map[string]interface{}{
		"status":  true,
		"message": "FCM Token Updated",
	}
	return c.JSON(http.StatusOK, responseBody)
}
