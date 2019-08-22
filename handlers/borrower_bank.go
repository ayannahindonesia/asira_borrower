package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func BorrowerBankService(c echo.Context) error {
	defer c.Request().Body.Close()

	bank_service_id, _ := strconv.Atoi(c.Param("bank_service_id"))

	bankService := models.BankService{}
	result, err := bankService.FindbyID(bank_service_id)
	if err != nil {
		return returnInvalidResponse(http.StatusNotFound, err, fmt.Sprintf("layanan %v tidak ditemukan", bank_service_id))
	}

	return c.JSON(http.StatusOK, result)
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
