package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/asira-ayannah/basemodel"

	"github.com/labstack/echo"
)

type BanksOfAgentSelect struct {
	models.Agent
	BankName string `json:"bank_name"`
}

type BankList struct {
	ID       int64  `json:"id"`
	BankName string `json:"bank_name"`
}

func AgentAllBank(c echo.Context) error {
	defer c.Request().Body.Close()
	var (
		agentBanks []BanksOfAgentSelect
		banklist   []BankList
		totalRows  int
		i          int
	)
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

	//get IDs of bank dari array pertama biar lebih aman
	bankIDs := agentBanks[0].Banks

	for _, val := range agentBanks {
		newBankList := BankList{
			ID:       bankIDs[i],
			BankName: val.BankName,
		}
		banklist = append(banklist, newBankList)
		i++
	}
	//create custom response
	tempDB := db
	tempDB.Count(&totalRows)
	result := basemodel.PagedFindResult{
		TotalData: totalRows,
		Data:      banklist,
	}
	return c.JSON(http.StatusOK, result)
}
