package handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//FUTURE: NotificationsGetByTopic
//NOTE:  handler for get notification data by recipient_id (custom format + borrower id; i.e. : borrower-%d)
func NotificationsGet(c echo.Context) error {
	defer c.Request().Body.Close()

	type Filter struct {
		RecipientID string `json:"recipient_id"`
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	notification := models.Notification{}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	//NOTE: borrower user default formater
	recipient_id := fmt.Sprintf("borrower-%d", borrowerID)

	//search
	result, err := notification.PagedFilterSearch(page, rows, orderby, sort, &Filter{
		RecipientID: recipient_id,
	})
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, "Notifikasi tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
