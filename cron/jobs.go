package cron

import (
	"log"
)

// AutoLoanDisburseConfirm confirms loan disburse status
func SendNotifications() func() {
	return func() {
		type Response struct {
			ID              uint64  `json:"id"`
			LoanPayment     float64 `json:"loan_payment"`
			InterestPayment float64 `json:"interest_payment"`
		}
		var responses []Response

		//get installments 3 day before due
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
			type LoanPaymentStatus struct {
				ID             uint64 `json:"id"`
				PaymentStatus  string `json:"payment_status"`
				DisburseStatus string `json:"disburse_status"`
				BorrowerID     uint64 `json:"borrower_id"`
				FCMToken       string `json:"fcm_token"`
			}

			var loanStatus LoanPaymentStatus
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
			if loanStatus.PaymentStatus == "processing" {
				sendRemainderNotif(loanStatus.BorrowerID, loanStatus.FCMToken)
			}
			log.Printf("SendNotifications cron executed. loanStatus : %+v ", loanStatus)
		}

		log.Printf("SendNotifications cron executed. response : %+v ", responses)
	}
}

func sendRemainderNotif(borrowerID uint64, fcmToken string) {
	log.Printf("send notif  : borrower = %+v ; token = %v ", borrowerID, fcmToken)

	// responseBody, err := asira.App.Messaging.SendNotificationByToken(title, message, mapData, user.FCMToken, recipientID)
}
