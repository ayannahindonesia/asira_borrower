package handlers

import (
	"asira_borrower/models"
	"fmt"
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

	type Filter struct {
		BankID string `json:"bank_id"`
		// NameOR []string `json:"name" condition:"OR"`
	}

	bankService := models.BankService{}
	result, err := bankService.PagedFindFilter(0, 0, []string{}, []string{}, &Filter{
		BankID: fmt.Sprintf("%v", borrowerModel.Bank.Int64),
	})

	return c.JSON(http.StatusOK, result)
}

func BorrowerBankServiceDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	bServices := models.BankService{}

	serviceID, _ := strconv.Atoi(c.Param("service_id"))
	err := bServices.FindbyID(serviceID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, bServices)
}
