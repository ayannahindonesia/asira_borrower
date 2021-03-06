package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//BorrowerResponse hold custom response
type BorrowerResponse struct {
	models.Borrower
	NthLoans   int    `json:"nth_loans" gorm:"-"`
	LoanStatus string `json:"loan_status"`
	BankName   string `json:"bank_name"`
}

//AgentAllBorrower get borrower list owned by agent
func AgentAllBorrower(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentAllBorrower"

	type Filter struct {
		AgentReferral sql.NullInt64 `json:"agent_referral"`
		Bank          sql.NullInt64 `json:"bank"`
	}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, err := strconv.ParseUint(claims["jti"].(string), 10, 64)
	var agent models.Agent
	err = agent.FindbyID(agentID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:    "not valid agent",
			NLOGERR:    err,
			"agent_id": agentID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	// pagination parameters

	db := asira.App.DB
	var (
		totalRows int
		offset    int
		rows      int
		page      int
		lastPage  int
		borrowers []BorrowerResponse
	)

	// pagination parameters
	rows, err = strconv.Atoi(c.QueryParam("rows"))
	page, err = strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	//params bank_id
	bankID, _ := strconv.ParseInt(c.QueryParam("bank_id"), 10, 64)

	//set rows, page and offset
	if rows > 0 {
		if page <= 0 {
			page = 1
		}
		offset = (page * rows) - rows
	}

	//query banks from agent's borrowers

	LoanStatusQuery := fmt.Sprintf("CASE WHEN (SELECT COUNT(id) FROM loans l WHERE l.borrower = borrowers.id AND status IN ('%s', '%s') AND (due_date IS NULL OR due_date = '0001-01-01 00:00:00+00' OR NOW() < l.due_date + make_interval(days => 1))) > 0 THEN '%s' ELSE '%s' END", "approved", "processing", "active", "inactive")
	//due_date IS NULL OR due_date = '0001-01-01 00:00:00+00' OR

	//filters
	db = db.Table("borrowers").
		Select("borrowers.*, "+LoanStatusQuery+" as loan_status, (SELECT COUNT(id) FROM loans l WHERE l.borrower = borrowers.id AND l.status = ?) as nth_loans, (SELECT b.name FROM banks b WHERE b.id = borrowers.bank) as bank_name", "approved").
		Where("borrowers.agent_referral = ?", agentID)

	//show borrowers for specific bankID
	if bankID > 0 {
		db = db.Where("borrowers.bank = ?", bankID)
	}

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
	countDB.Where("borrowers.deleted_at IS NULL").Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}
	err = db.Find(&borrowers).Error
	if err != nil {

		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:    "empty agent's borrower list",
			NLOGERR:    err,
			"agent_id": agentID}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusInternalServerError, err, "data agent's borrowers tidak ditemukan")
	}

	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        borrowers,
	}

	return c.JSON(http.StatusOK, result)
}
