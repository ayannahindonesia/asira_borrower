package cron

import (
	"asira_borrower/custommodule"
	"fmt"
	"log"

	"encoding/json"
)

//Installment installment response
type Installment struct {
	ID              uint64  `json:"id"`
	LoanPayment     float64 `json:"loan_payment"`
	InterestPayment float64 `json:"interest_payment"`
}

//LoanPayment loan response
type LoanPayment struct {
	ID             uint64 `json:"id"`
	PaymentStatus  string `json:"payment_status"`
	DisburseStatus string `json:"disburse_status"`
	BorrowerID     uint64 `json:"borrower_id"`
	FCMToken       string `json:"fcm_token"`
}

//SendInstallmentNotifications confirms loan disburse status
func SendInstallmentNotifications() func() {
	return func() {
		var responses []Installment

		//get installments 3 day before due date
		err := DB.Table("installments").
			Where("due_date BETWEEN DATE(now()+make_interval(days => 3)) AND DATE(now()+make_interval(days => 4))").
			Find(&responses).Error
		if err != nil {
			log.Printf("SendNotifications cron executed. error : %v", err)
			return
		}

		//cek for loan payment_status for each installements
		for _, res := range responses {
			// log.Printf("%+v\n", res)
			var loanStatus LoanPayment
			err = DB.Table("loans").
				Select("loans.id, b.id AS borrower_id, payment_status, disburse_status, fcm_token").
				Joins("INNER JOIN borrowers b ON b.id = loans.borrower").
				Joins("INNER JOIN users u ON u.borrower = b.id").
				Where("? IN (SELECT UNNEST(loans.installment_details) )", res.ID).
				Find(&loanStatus).Error
			if err != nil {
				log.Printf("installment dont have valid parent (loan) id")
				continue
			}

			//cek current status loan (disburse_status == confirmed)
			if loanStatus.PaymentStatus == "processing" && loanStatus.DisburseStatus == "confirmed" {
				_ = sendRemainderNotif(loanStatus, res)
			}
			log.Printf("SendNotifications cron executed. loanStatus : %+v ", loanStatus)
		}

		log.Printf("SendNotifications cron executed. response : %+v ", responses)
	}
}

func sendRemainderNotif(loan LoanPayment, installment Installment) error {
	log.Printf("send notif  : borrower = %+v ; token = %v ", loan.BorrowerID, loan.FCMToken)

	recipientID := fmt.Sprintf("borrower-%d", loan.BorrowerID)
	title := fmt.Sprintf("Cicilan Pembayaran Pinjaman %d", loan.ID)

	//TODO: pinalty cek
	message := fmt.Sprintf("Cicilan anda akan masuk masa jatuh tempo dalam 3 hari, silahkan lakukan pembayaran sebesar Rp.%0.0f", installment.LoanPayment+installment.InterestPayment)

	//send push notification
	responseBody, err := custommodule.MessagingStatic.SendNotificationByToken(title, message, nil, loan.FCMToken, recipientID)

	// if error create error info in notifications table
	if err != nil {
		type ErrorResponse struct {
			Details string `json:"details"`
			Message string `json:"message"`
		}
		var errorResponse ErrorResponse

		//parse error response
		err = json.Unmarshal(responseBody, &errorResponse)
		if err != nil {
			log.Printf(err.Error())
			return err
		}
		return err
	}
	return nil
}
