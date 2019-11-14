package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.com/asira-ayannah/basemodel"
)

//FUTURE: NotificationsGetByTopic
//NOTE:  handler for get notification data by recipient_id (custom format + borrower id; i.e. : borrower-%d)
func NotificationsGet(c echo.Context) error {
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

	//NOTE: borrower user default formater
	recipient_id := fmt.Sprintf("borrower-%d", borrowerID)

	//get notification list by recipient_id
	response, err := asira.App.Messaging.GetNotificationByRecipientID(recipient_id, c)
	if err != nil {
		return err //returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending notification")
	}

	//parse result from microservice messaging
	var parseResponse basemodel.PagedFindResult
	json.Unmarshal([]byte(response), &parseResponse)
	return c.JSON(http.StatusOK, parseResponse)
}
