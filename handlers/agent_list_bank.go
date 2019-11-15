package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

type BanksOfAgentSelect struct {
	models.Agent
	BankName string `json:"bank_name"`
}

func AgentProfile(c echo.Context) error {
	defer c.Request().Body.Close()
	var agentBanks []BanksOfAgentSelect

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.Atoi(claims["jti"].(string))

	//custom query
	db := asira.App.DB
	db = db.Table("agents a").
		Select("a.*, b.name as bank_name").
		Joins("INNER JOIN banks b ON b.id = ANY(a.banks)").
		Where("a.id = ?", agentID)
	//query
	err = db.Find(&agentBanks).Error
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	return c.JSON(http.StatusOK, agentBanks)
}
