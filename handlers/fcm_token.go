package handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

//FCMTokenUpdate update FCMToken for current user device
func FCMTokenUpdate(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "FCMTokenUpdate"

	var payload models.User

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userBorrower := models.User{}

	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := userBorrower.FindbyBorrowerID(borrowerID)
	if err != nil {
		NLog("warning", LogTag, fmt.Sprintf("unauthorized : %v", borrowerID), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

	payloadRules := govalidator.MapData{
		"fcm_token": []string{},
	}

	validate := validateRequestPayload(c, payloadRules, &payload)
	if validate != nil {
		NLog("warning", LogTag, fmt.Sprintf("error validation : %v", validate), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	//update FCMToken
	userBorrower.FCMToken = payload.FCMToken
	err = userBorrower.Save()
	if err != nil {
		NLog("warning", LogTag, fmt.Sprintf("error update user borrower FCMToken : %v", err), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "error saving Password")
	}

	responseBody := map[string]interface{}{
		"status":  true,
		"message": "FCM Token Updated",
	}

	NLog("event", LogTag, fmt.Sprintf("success updating FCMToken : %v", userBorrower.FCMToken), c.Get("user").(*jwt.Token), "", false, "borrower")

	return c.JSON(http.StatusOK, responseBody)
}
