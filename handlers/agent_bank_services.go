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

func AgentBankService(c echo.Context) error {

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
	agentID, _ := strconv.Atoi(claims["jti"].(string))
	err := agentModel.FindbyID(agentID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun agen tidak ditemukan")
	}

	//check service by bank_id
	bankID, _ := strconv.ParseInt(c.Param("bank_id"), 10, 64)

	//check bank exist in Agent.Banks; manual looping for performance
	exist := false
	for _, val := range agentModel.Banks {
		if val == bankID {
			exist = true
			break
		}
	}
	if exist == false {
		returnInvalidResponse(http.StatusForbidden, err, "Bank ID tidak terdaftar untuk agen")
	}
	fmt.Println("bankID=", bankID, "; agentID=", agentID)
	//query result
	db := asira.App.DB
	var results []models.Service
	var count int
	err = db.Table("services").
		Select("*").
		Where("id IN (SELECT UNNEST(services) FROM banks WHERE id = ?)", bankID).Find(&results).Count(&count).Error

	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Product Tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, &Result{TotalData: count, Data: results})
}

func AgentBankServiceDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	bServices := models.Service{}

	serviceID, _ := strconv.Atoi(c.Param("service_id"))
	err := bServices.FindbyID(serviceID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, bServices)
}
