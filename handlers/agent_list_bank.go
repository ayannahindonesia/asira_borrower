package handlers

import (
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func AgentAllBank(c echo.Context) error {
	defer c.Request().Body.Close()
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
		return returnInvalidResponse(http.StatusInternalServerError, err, "data agent banks tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
