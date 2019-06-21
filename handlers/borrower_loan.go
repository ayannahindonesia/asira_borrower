package handlers

import (
	"asira/models"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func BorrowerLoanApply(c echo.Context) error {
	defer c.Request().Body.Close()

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loan.Owner = sql.NullInt64{Int64: int64(borrowerID), Valid: true}

	payloadRules := govalidator.MapData{
		"loan_amount":       []string{"required"},
		"installment":       []string{"required"},
		"loan_intention":    []string{"required"},
		"intention_details": []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &loan)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	newLoan, err := loan.Create()
	if err != nil {
		log.Printf("apply : %v", loan)
		return returnInvalidResponse(http.StatusInternalServerError, err, "create new loan failed")
	}

	return c.JSON(http.StatusCreated, newLoan)
}

func BorrowerLoanGet(c echo.Context) error {
	defer c.Request().Body.Close()

	loan := models.Loan{}

	// pagination parameters
	rows, err := strconv.Atoi(c.QueryParam("rows"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	orderby := c.QueryParam("orderby")
	sort := c.QueryParam("sort")

	// filters
	status := c.QueryParam("status")

	payloadRules := govalidator.MapData{
		"rows": []string{"numeric"},
		"page": []string{"numeric"},
	}

	validate := validateRequestQuery(c, payloadRules)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	result, err := loan.PagedSearch(page, rows, orderby, sort, &models.LoanSearchFilter{
		Status: status,
	})

	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, err, "query result error")
	}

	return c.JSON(http.StatusOK, result)
}
