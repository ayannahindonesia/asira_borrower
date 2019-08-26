package handlers

import (
	"asira_borrower/models"
	"encoding/json"
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

	// var service []int
	// jMarshal, _ := json.Marshal(bankBorrower.Services)
	// if err := json.Unmarshal(jMarshal, &service); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(service)
	// bankService := models.BankService{}
	// bServices := make([]interface{}, len(service))
	// for i := range service {
	// 	data, err := bankService.FindbyID(service[i])
	// 	if err != nil {
	// 		continue
	// 	}
	// 	bServices[i] = data
	// }

	return c.JSON(http.StatusOK, bankBorrower.Services)
}

func BorrowerBankServiceDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	bServices := models.BankService{}

	serviceID, _ := strconv.Atoi(c.Param("service_id"))
	service, err := bServices.FindbyID(serviceID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, service)
}

func BorrowerBankProduct(c echo.Context) error {
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

	var product []int
	jMarshal, _ := json.Marshal(bankBorrower.Products)
	if err := json.Unmarshal(jMarshal, &product); err != nil {
		log.Fatal(err)
	}

	bankProduct := models.ServiceProduct{}
	bProduct := make([]interface{}, len(product))
	for i := range product {
		data, err := bankProduct.FindbyID(product[i])
		if err != nil {
			continue
		}
		bProduct[i] = data
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total_data": len(bProduct),
		"data":       bProduct,
	})
}

func BorrowerBankProductDetails(c echo.Context) error {
	defer c.Request().Body.Close()
	sProduct := models.ServiceProduct{}

	productID, _ := strconv.Atoi(c.Param("product_id"))
	product, err := sProduct.FindbyID(productID)
	if err != nil {
		return returnInvalidResponse(http.StatusForbidden, err, "Service Product Tidak Ditemukan")
	}
	return c.JSON(http.StatusOK, product)
}
