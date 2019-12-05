package handlers

import (
	"asira_borrower/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func AgentAllBorrower(c echo.Context) error {
	defer c.Request().Body.Close()
	var (
		rows int
		page int
	)
	type Filter struct {
		AgentID sql.NullInt64 `json:"agent_id"`
		Bank    sql.NullInt64 `json:"bank"`
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, err := strconv.ParseInt(claims["jti"].(string), 10, 64)
	var agent models.Agent
	err = agent.FindbyID(int(agentID))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	// pagination parameters
	rows, err = strconv.Atoi(c.QueryParam("rows"))
	page, err = strconv.Atoi(c.QueryParam("page"))
	order := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	//params bank_id
	bankID, _ := strconv.ParseInt(c.Param("bank_id"), 10, 64)

	//query banks from agent's borrowers
	var AgentBorrower models.Borrower
	result, err := AgentBorrower.PagedFilterSearch(page, rows, order, sort, &Filter{
		AgentID: sql.NullInt64{
			Int64: agentID,
			Valid: true,
		},
		Bank: sql.NullInt64{
			Int64: bankID,
			Valid: true,
		},
	})
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "data agent's borrowers tidak ditemukan")
	}

	return c.JSON(http.StatusOK, result)
}
