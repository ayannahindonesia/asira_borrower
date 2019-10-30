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

	loan.Owner = sql.NullInt64{Int64: int64(borrowerID), Valid: true}

	payloadRules := govalidator.MapData{
		"loan_amount":       []string{"required"},
		"installment":       []string{"required"},
		"loan_intention":    []string{"required", "loan_purposes"},
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
	err = LoanSimulateApproveReject(loan, simulate)
	if err != nil {
		return returnInvalidResponse(http.StatusInternalServerError, loan, err.Error()) //"failed to simulate loan approve / reject : ",
	}
	// END ****** simulate loan approve / reject *******
	return c.JSON(http.StatusCreated, loan)
}

func LoanSimulateApproveReject(loan models.Loan, status string) error {

	loan.Status = status
	loan.DisburseDate = time.Now()

	//update status
	err = loan.SaveNoKafka()
	if err != nil {
		return err
	}
	jsonReq, _ := json.Marshal(loan)
	//TODO: send notification
	err := asira.App.Messaging.SendNotificationByToken("testing", string(jsonReq), "cEh41s_l_t4:APA91bGaE1OLrCN0P3myiSslwtddtmZMDj4uy_0YbJJ3qvt_N_f81HdxJL5juuuud18OW3zfKZqLDMbn83O1EoBBhGHvJMKupupb5CUsSaWc9A4b6bItmDEctwZ3F-5ENoJfHPZP4NMn")
	if err != nil {
		return err //returnInvalidResponse(http.StatusUnprocessableEntity, err, "failed sending notification")
	}
	//err = models.KafkaSubmitModelLoopback(loan, "loan")
	return nil
}
