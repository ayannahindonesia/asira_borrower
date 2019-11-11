package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

func getBorrowerId(c echo.Context) (int, error) {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrower := models.Borrower{}
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrower.FindbyID(borrowerID)
	if err != nil {
		return 0, err
	}
	return borrowerID, nil
}

//FUTURE: NotificationsGetByTopic
//TODO: handler for get notification data by borrower id (NotificationsGetBySenderID)
//NOTE:  handler for get notification data by FCM Token
func NotificationsGetByToken(c echo.Context) error {
	defer c.Request().Body.Close()
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrower := models.Borrower{}
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrower.FindbyID(borrowerID)
	// borrowerID, err := getBorrowerId(c)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}

	response, err := asira.App.Messaging.GetNotificationByToken(borrower.FCMToken)
	if err != nil {
		return err //returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending notification")
	}

	return c.JSON(http.StatusOK, response)
}
