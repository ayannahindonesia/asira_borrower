package handlers

import (
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//AgentAllBank get bank list owned by agent
func AgentAllBank(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentAllBank"

	var (
		rows int
		page int
	)
	type Filter struct {
		ID []int64 `json:"id"`
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, err := strconv.ParseUint(claims["jti"].(string), 10, 64)
	var agent models.Agent
	err = agent.FindbyID(agentID)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("not valid agent : %v agent ID", err, agentID), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	// pagination parameters
	rows, err = strconv.Atoi(c.QueryParam("rows"))
	page, err = strconv.Atoi(c.QueryParam("page"))
	order := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	//query banks from agent's banks
	var banks models.Bank
	result, err := banks.PagedFilterSearch(page, rows, order, sort, &Filter{
		ID: []int64(agent.Banks),
	})
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("not found agent's banks : %v", err), c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusInternalServerError, err, "data agent banks tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
