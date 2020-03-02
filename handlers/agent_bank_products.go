package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//AgentBankProduct get products in agent.banks
func AgentBankProduct(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "AgentBankProduct"

	type Result struct {
		TotalData int              `json:"total_data"`
		Data      []models.Product `json:"data"`
	}

	var results []models.Product
	var count int
	db := asira.App.DB

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	agentID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	agentModels := models.Agent{}
	err := agentModels.FindbyID(agentID)
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "not valid agent",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusForbidden, err, "agent tidak valid")
	}

	db = db.Table("products").
		Select("products.*").
		Joins("INNER JOIN services s ON s.id = products.service_id").
		Joins("INNER JOIN banks b ON s.id IN (SELECT UNNEST(b.services)) AND products.id IN (SELECT UNNEST(b.products) )").
		Joins("INNER JOIN agents ag ON b.id IN (SELECT UNNEST(ag.banks))").
		Where("ag.id = ?", agentID)

	//query tambahan jika parameter terdefinisi
	if serviceID := c.QueryParam("service_id"); len(serviceID) > 0 {
		db = db.Where("s.id = ?", serviceID)
	}
	if productID := c.QueryParam("product_id"); len(productID) > 0 {
		db = db.Where("products.id = ?", productID)
	}

	//harus di group by krn dr ag.banks lebih dari 1
	db = db.Group("products.id")

	err = db.Find(&results).Count(&count).Error
	if err != nil || count == 0 {
		NLog("warning", LogTag, map[string]interface{}{
			NLOGMSG:   "empty products list",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", false, "agent")

		return returnInvalidResponse(http.StatusNotFound, err, "Service Product Tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, &Result{TotalData: count, Data: results})
}
