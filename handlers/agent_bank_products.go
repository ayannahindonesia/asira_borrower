package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func AgentBankProduct(c echo.Context) error {
	defer c.Request().Body.Close()

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
		return returnInvalidResponse(http.StatusForbidden, err, "agent tidak valid")
	}

	db = db.Table("banks").
		Select("p.*").
		Joins("INNER JOIN agents ag ON banks.id IN (SELECT UNNEST(ag.banks))").
		Joins("INNER JOIN services s ON s.id IN (SELECT UNNEST(banks.services))").
		Joins("INNER JOIN products p ON p.service_id = s.id ").
		Where("ag.id = ?", agentID).
		Where(generateDeleteCheck("p"))

	//query tambahan jika parameter terdefinisi
	if serviceID := c.QueryParam("service_id"); len(serviceID) > 0 {
		db = db.Where("s.id = ?", serviceID)
	}
	if productID := c.QueryParam("product_id"); len(productID) > 0 {
		db = db.Where("p.id = ?", productID)
	}

	//harus di group by krn dr ag.banks lebih dari 1
	db = db.Group("p.id")

	err = db.Find(&results).Count(&count).Error
	if err != nil || count == 0 {
		return returnInvalidResponse(http.StatusNotFound, err, "Service Product Tidak Ditemukan")
	}

	return c.JSON(http.StatusOK, &Result{TotalData: count, Data: results})
}
