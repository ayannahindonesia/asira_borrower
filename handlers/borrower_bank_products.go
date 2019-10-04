package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func BorrowerBankProduct(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	bank := models.Bank{}
	bank.FindbyID(int(borrowerModel.Bank.Int64))

	db := asira.App.DB
	var results []models.BankProduct
	var count int

	err = db.Table("bank_products p").
		Select("p").
		Joins("INNER JOIN bank_services s ON s.id = p.bank_service_id").
		Joins("INNER JOIN banks b ON b.id = s.bank_id").
		Joins("INNER JOIN borrowers bo ON bo.bank = b.id").
		Where("bo.id = ?", borrowerID).Scan(&results).Count(&count).Error

	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Product Tidak Ditemukan")
	}

	type Result struct {
		TotalData int                  `json:"total_data"`
		Data      []models.BankProduct `json:"data"`
	}

	return c.JSON(http.StatusOK, &Result{TotalData: count, Data: results})
}

func BorrowerBankProductDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	bankProduct := models.BankProduct{}

	productID, _ := strconv.Atoi(c.Param("product_id"))
	err := bankProduct.FindbyID(productID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Product Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, bankProduct)
}
