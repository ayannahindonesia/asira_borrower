package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/lib/pq"
)

//AgentBankService get services in agent.banks
func AgentBankService(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentBankService"

	type Filter struct {
		Banks pq.Int64Array `json:"banks"`
	}

	type Result struct {
		TotalData int              `json:"total_data"`
		Data      []models.Service `json:"data"`
	}

	//cek agent id
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentModel := models.Agent{}
	agentID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := agentModel.FindbyID(agentID)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("not valid agent : %v agent ID : %v", err, agentID), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun agen tidak ditemukan")
	}

	//get from QueryParam
	bankID, _ := strconv.ParseInt(c.QueryParam("bank_id"), 10, 64)
	serviceID, _ := strconv.Atoi(c.QueryParam("service_id"))

	//check bank exist in Agent.Banks; manual looping for performance
	if isInArrayInt64(bankID, []int64(agentModel.Banks)) == false {
		NLog("warning", LogTag, fmt.Sprintf("not valid bank ID : %v", bankID), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusForbidden, err, "Bank ID tidak terdaftar untuk agen")
	}

	//query result serviceID
	db := asira.App.DB
	var results []models.Service
	var count int

	//build query
	objDB := db.Table("services").
		Select("*").
		Where("services.id IN (SELECT UNNEST(services) FROM banks b WHERE b.id = ?)", bankID)

	//bila serviceID di set berarti mengarah ke detail ID
	if serviceID > 0 {
		// bankServices := models.Service{}
		objDB = objDB.Where("services.id = ?", serviceID)
	}

	err = objDB.Find(&results).Count(&count).Error

	if err != nil || count == 0 {
		NLog("error", LogTag, fmt.Sprintf("service not found : %v count : %v", err, count), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusNotFound, err, "Service Tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, &Result{TotalData: count, Data: results})
}
