package handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

func getBorrower (c echo.Context) (models.Borrower){
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrower := models.Borrower{}
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrower.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "unauthorized")
	}
	return borrower
}

//FUTURE: NotificationsGetByTopic
//NOTE: handler for get notification data by borrower id
func NotificationsGetBySenderID(c echo.Context) error {
	defer c.Request().Body.Close()

	borrower := getBorrower (c) 

	err = model.RefreshToken()
	if err != nil {
		return err
	}
	if firebase_token == "" {
		firebase_token = model.Endpoints.PushNotification
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"title":          title,
		"message_body":   message_body,
		"firebase_token": firebase_token,
		"data":           map_data,
	})

	request, _ := http.NewRequest("POST", model.URL+model.Endpoints.PushNotification, bytes.NewBuffer(payload))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", model.Token))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
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
