package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//BorrowerBankProduct get list product
func BorrowerBankProduct(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	LogTag := "BorrowerBankProduct"

	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	db := asira.App.DB
	var results []models.Product
	var count int

	db = db.Table("products").
		Select("products.*").
		Joins("INNER JOIN services s ON s.id = products.service_id").
		Joins("INNER JOIN banks bnk ON s.id IN (SELECT UNNEST(bnk.services)) AND products.id IN (SELECT UNNEST(bnk.products))").
		Joins("INNER JOIN borrowers bo ON bo.bank = bnk.id").
		Where("bo.id = ?", borrowerID)

	if serviceID := c.QueryParam("service_id"); len(serviceID) > 0 {
		db = db.Where("s.id = ?", serviceID)
	}

	err = db.Find(&results).Count(&count).Error
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "Product not found",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Product Tidak Ditemukan")
	}

	type Result struct {
		TotalData int              `json:"total_data"`
		Data      []models.Product `json:"data"`
	}

	return c.JSON(http.StatusOK, &Result{TotalData: count, Data: results})
}

//BorrowerBankProductDetails get details
func BorrowerBankProductDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	LogTag := "BorrowerBankProductDetails"

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	results := models.Product{}
	productID, _ := strconv.ParseUint(c.Param("product_id"), 10, 64)

	db := asira.App.DB
	db = db.Table("products").
		Select("products.*").
		Joins("INNER JOIN services s ON s.id = products.service_id").
		Joins("INNER JOIN banks bnk ON s.id IN (SELECT UNNEST(bnk.services)) AND products.id IN (SELECT UNNEST(bnk.products))").
		Joins("INNER JOIN borrowers bo ON bo.bank = bnk.id").
		Where("bo.id = ?", borrowerID).
		Where("products.id = ?", productID)

	err = db.Find(&results).Error
	if err != nil {
		NLog("error", LogTag, map[string]interface{}{
			NLOGMSG:   "Product not found",
			NLOGERR:   err,
			NLOGQUERY: asira.App.DB.QueryExpr()}, c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Product Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, results)
}
