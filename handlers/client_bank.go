package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/lib/pq"
)

type BankResponse struct {
	models.Bank
	ServiceName pq.StringArray `json:"service_name"`
}

//ClientBankServices get service
func ClientBanks(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "ClientBanks"

	var (
		totalRows int
		offset    int
		rows      int
		page      int
		lastPage  int
		banks     []BankResponse
	)
	db := asira.App.DB

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err = strconv.Atoi(c.QueryParam("page"))
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")
	//get bank id
	bankID, _ := strconv.ParseUint(c.QueryParam("bank_id"), 10, 64)

	//set rows, page and offset
	if rows > 0 {
		if page <= 0 {
			page = 1
		}
		offset = (page * rows) - rows
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

	//build query
	db = db.Table("banks").
		Select("*, (SELECT ARRAY_AGG(s.name) FROM services s WHERE s.id IN (SELECT UNNEST(banks.services) ) ) as service_name")

	if bankID > 0 {
		db = db.Where("banks.id = ?", bankID)
	}
	// countDB := db
	// countDB.Where("services.deleted_at IS NULL").Count(&totalRows)

	if rows > 0 {
		db = db.Limit(rows).Offset(offset)
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))
	}

	err = db.Find(&banks).Error
	if err != nil {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "empty banks list",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, err, "empty banks list")
	}

	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        banks,
	}

	return c.JSON(http.StatusOK, result)
}

//ClientBankbyID get bank detail
func ClientBankbyID(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "ClientBankbyID"

	bank := models.Bank{}
	bankID, _ := strconv.ParseUint(c.Param("bank_id"), 10, 64)
	err := bank.FindbyID(bankID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "query result error",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", true, "")

		return returnInvalidResponse(http.StatusInternalServerError, err, "bank tidak ditemukan")
	}

	return c.JSON(http.StatusOK, bank)
}
