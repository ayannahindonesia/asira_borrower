package handlers

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thedevsaddam/govalidator"
)

func LoanNotificationSimulate(c echo.Context) error {
	defer c.Request().Body.Close()
	var err error

	loan := models.Loan{}

	user := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	borrowerID, _ := strconv.Atoi(claims["jti"].(string))

	loan.Borrower = sql.NullInt64{Int64: int64(borrowerID), Valid: true}

	payloadRules := govalidator.MapData{
		"loan_amount":       []string{"required"},
		"installment":       []string{"required"},
		"loan_intention":    []string{"required"},
		"intention_details": []string{"required"},
		"product":           []string{"required"},
	}

	validate := validateRequestPayload(c, payloadRules, &loan)
	if validate != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "validation error")
	}

	err = validateLoansProduct(loan)
	if err != nil {
		return returnInvalidResponse(http.StatusUnprocessableEntity, err, "validation error")
	}
	err = loan.Create()
	if err != nil {
		log.Printf("apply : %v", loan)
		return returnInvalidResponse(http.StatusInternalServerError, err, "Gagal membuat Loan")
	}
	log.Printf("apply : %v", loan)

	// ****** simulate loan approve / reject *******
	simulate := c.QueryParam("simulate")
	if simulate != "approve" && simulate != "reject" {
		return returnInvalidResponse(http.StatusUnprocessableEntity, validate, "simulate = [ \"approve\" | \"reject\" ]")
	}
	tokenparam := c.QueryParam("token")
	err = LoanSimulateApproveReject(&loan, simulate, tokenparam)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, loan, err.Error()) //"failed to simulate loan approve / reject : ",
	}
	// END ****** simulate loan approve / reject *******
	return c.JSON(http.StatusCreated, loan)
}

func LoanSimulateApproveReject(loan *models.Loan, status string, token string) error {

	loan.Status = status
	loan.DisburseDate = time.Now()

	//update status
	err = loan.SaveNoKafka()
	if err != nil {
		return err
	}
	jsonReq, _ := json.Marshal(loan)
	//MAYBEDO: general type (string, int, float) to map[string]string
	err := asira.App.Messaging.SendNotificationByToken("testing", string(jsonReq), nil, token, "borrower-simulate")
	if err != nil {
		return err //returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending notification")
	}
	//err = models.KafkaSubmitModelLoopback(loan, "loan")
	return nil
}
