package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

//BorrowerBankService get bank's service list
func BorrowerBankService(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	LogTag := "BorrowerBankService"

	borrowerModel := models.Borrower{}

	borrowerID, _ := strconv.ParseUint(claims["jti"].(string), 10, 64)
	err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("error borrower ID : %v", err), c.Get("user").(*jwt.Token), "", false, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	db := asira.App.DB
	var results []models.Service
	var count int

	err = db.Table("services").
		Select("*").
		Where("id IN (SELECT UNNEST(services) FROM banks WHERE id = ?)", borrowerModel.Bank.Int64).Find(&results).Count(&count).Error

	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("Service not found"), c.Get("user").(*jwt.Token), "", true, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Service Tidak Ditemukan")
	}

	type Result struct {
		TotalData int              `json:"total_data"`
		Data      []models.Service `json:"data"`
	}

	return c.JSON(http.StatusOK, &Result{TotalData: count, Data: results})
}

//BorrowerBankServiceDetails get service detail
func BorrowerBankServiceDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	bServices := models.Service{}

	LogTag := "BorrowerBankServiceDetails"

	serviceID, _ := strconv.ParseUint(c.Param("service_id"), 10, 64)
	err := bServices.FindbyID(serviceID)
	if err != nil {
		NLog("error", LogTag, fmt.Sprintf("Service detail not found"), c.Get("user").(*jwt.Token), "", true, "borrower")

		return returnInvalidResponse(http.StatusForbidden, err, "Service Detail Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, bServices)
}
