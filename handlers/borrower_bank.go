package handlers

import (
	"asira_borrower/models"
	"encoding/json"
	"fmt"
	"log"
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
	borrower, err := borrowerModel.FindbyID(borrowerID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Akun tidak ditemukan")
	}

	bank := models.Bank{}
	bankBorrower, _ := bank.FindbyID(int(borrower.Bank.Int64))

	var service []int

	if err := json.Unmarshal(bankBorrower.Products, &service); err != nil {
		log.Fatal(err)
	}

	bankService := models.BankService{}
	bServices := make([]interface{}, len(service))
	for i, s := range service {
		bServices[i], err = bankService.FindbyID(service[i])
	}

	return c.JSON(http.StatusOK, bServices)
}

func BorrowerBankProduct(c echo.Context) error {
	defer c.Request().Body.Close()

	product_id, _ := strconv.Atoi(c.Param("product_id"))

	serviceProduct := models.ServiceProduct{}
	result, err := serviceProduct.FindbyID(product_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("product %v tidak ditemukan", serviceProduct))
	}

	return c.JSON(http.StatusOK, result)
}
