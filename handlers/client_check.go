package handlers

import (
	"asira/asira"
	"asira/models"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func CheckUniqueData(c echo.Context) error {
	defer c.Request().Body.Close()

	borrower := models.Borrower{}

	payloadRules := govalidator.MapData{
		"idcard_number": []string{"unique:borrowers,idcard_number"},
		"email":         []string{"email", "unique:borrowers,email"},
		"phone":         []string{"unique:borrowers,phone"},
	}

	validate := validateRequestPayload(c, payloadRules, &borrower)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	return c.JSON(http.StatusOK, newBorrower)
}
