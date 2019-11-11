package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.com/asira-ayannah/basemodel"
)

//FUTURE: NotificationsGetByTopic
//TODO: handler for get notification data by QueryParam AND borrower id (NotificationsGetBySenderID)
//NOTE:  handler for get notification data by QueryParam
func NotificationsGetByToken(c echo.Context) error {
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

	//get data
	response, err := asira.App.Messaging.GetNotificationByToken(borrower.FCMToken, c)
	if err != nil {
		return err //returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending notification")
	}
	var parseResponse basemodel.PagedFindResult
	json.Unmarshal([]byte(response), &parseResponse)
	return c.JSON(http.StatusOK, parseResponse)
}
