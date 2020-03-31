package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"asira_borrower/modules"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//AgentAllBank get bank list owned by agent
func AgentAllBank(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentAllBank"

	var agent models.Agent
	var banks []BankResponse

	type Filter struct {
		BankID int64  `json:"id"`
		Name   string `json:"name"`
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, err := strconv.ParseUint(claims["jti"].(string), 10, 64)

	err = agent.FindbyID(agentID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:    "not valid agent",
			NLOGERR:    err,
			"agent_id": agentID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	//Extended Query
	QPaged := modules.QueryPaged{}
	QPaged.Init(c)
	// fmt.Println("QPaged = ", QPaged)

	//custom query
	db := asira.App.DB

	//build query

	db = db.Table("banks").
		Select("*, (SELECT ARRAY_AGG(s.name) FROM services s WHERE s.id IN (SELECT UNNEST(banks.services) ) ) as service_name").
		Where("banks.id IN (?)", []int64(agent.Banks))

	//generate filter, return db and error
	db, err = QPaged.GenerateFilters(db, Filter{}, "banks")
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			"message": "error listing services",
			"error":   err}, token, "", false, "agent")

		return returnInvalidResponse(http.StatusInternalServerError, err, "kesalahan dalam filters")
	}

	//execute anonymous function pass db and data pass by reference (services)
	err = QPaged.Exec(db, &banks, func(DB *gorm.DB, rows interface{}) error {
		//manual type casting :)
		err := DB.Find(rows.(*[]BankResponse)).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			"message": "error listing services",
			"error":   err}, token, "", false, "agent")

		return returnInvalidResponse(http.StatusInternalServerError, err, "Pencarian tidak ditemukan")
	}

	//get result format
	result := QPaged.GetPage(banks)

	return c.JSON(http.StatusOK, result)
}
