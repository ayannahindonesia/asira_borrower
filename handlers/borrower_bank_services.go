package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func BorrowerBankService(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.Atoi(claims["jti"].(string))
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	db := asira.App.DB
	var results []models.Service
	var count int

	err = db.Table("services").
		Select("*").
		Where("id IN (SELECT UNNEST(services) FROM banks WHERE id = ?)", borrowerModel.Bank.Int64).Find(&results).Count(&count).Error

	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Product Tidak Ditemukan")
	}

	type Result struct {
		TotalData int              `json:"total_data"`
		Data      []models.Service `json:"data"`
	}

	return c.JSON(http.StatusOK, &Result{TotalData: count, Data: results})
}

func BorrowerBankServiceDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	bServices := models.Service{}

	serviceID, _ := strconv.Atoi(c.Param("service_id"))
	err := bServices.FindbyID(serviceID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, bServices)
}
