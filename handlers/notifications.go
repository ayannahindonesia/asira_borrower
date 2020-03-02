package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//NotificationsGet  handler for get notification data by recipient_id (custom format + borrower id; i.e. : borrower-%d)
//FUTURE: NotificationsGetByTopic
func NotificationsGet(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "NotificationsGet"

	type Filter struct {
		RecipientID string `json:"recipient_id"`
		Title       string `json:"title"`
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	db := asira.App.DB
	var (
		totalRows     int
		offset        int
		rows          int
		page          int
		lastPage      int
		notifications []models.Notification
	)

	// pagination parameters
	rows, err = strconv.Atoi(c.QueryParam("rows"))
	page, err = strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	// pagination parameters
	if rows > 0 {
		if page <= 0 {
			page = 1
		}
		offset = (page * rows) - rows
	}

	//NOTE: borrower user default formater
	recipientID := fmt.Sprintf("borrower-%d", borrowerID)

	//filters
	db = db.Table("notifications").
		Select("*").
		Where("recipient_id = ?", recipientID).
		Where("title <> 'failed'")

	if len(orderby) > 0 {
		if len(sort) > 0 {
			for k, v := range orderby {
				q := v
				if len(sort) > k {
					value := sort[k]
					if strings.ToUpper(value) == "ASC" || strings.ToUpper(value) == "DESC" {
						q = v + " " + strings.ToUpper(value)
					}
				}
				db = db.Order(q)
			}
		}
	}

	countDB := db
	countDB.Where("notifications.deleted_at IS NULL").Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}

	err = db.Find(&notifications).Error
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "empty query notification",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", false, "borrower")
	}

	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        notifications,
	}
	return c.JSON(http.StatusOK, result)
}
