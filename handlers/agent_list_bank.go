package handlers

import (
	"asira_borrower/asira"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/asira-ayannah/basemodel"

	"github.com/labstack/echo"
)

type BankList struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func AgentAllBank(c echo.Context) error {
	defer c.Request().Body.Close()
	var (
		banklist  []BankList
		totalRows int
		offset    int
		rows      int
		page      int
	)
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, err := strconv.Atoi(claims["jti"].(string))
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	// pagination parameters
	rows, _ = strconv.Atoi(c.QueryParam("rows"))
	if rows > 0 {
		page, _ = strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
		offset = (page * rows) - rows
	}

	//custom query
	db := asira.App.DB
	db = db.Table("banks b").
		Select("b.id, b.name").
		Joins("LEFT JOIN agents a ON b.id = ANY(a.banks)").
		Where("a.id = ?", agentID)

	//set all Paged properties
	tempDB := db
	tempDB.Count(&totalRows)

	//querying
	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
	} else {
		rows = totalRows
	}
	err = db.Find(&banklist).Error
	if err != nil {
		fmt.Println(err)
	}

	lastPage := int(math.Ceil(float64(totalRows) / float64(rows)))

	//create custom response
	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        banklist,
	}
	return c.JSON(http.StatusOK, result)
}
