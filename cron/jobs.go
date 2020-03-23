package cron

import (
	"log"
)

// AutoLoanDisburseConfirm confirms loan disburse status
func SendNotifications() func() {
	return func() {
		type Response struct {
			ID              uint64 `json:"id"`
			InterestPayment uint64 `json:"interest_payment"`
		}
		var responses []Response
		err := DB.Table("installments").
			Where("extract('day' from date_trunc('day', NOW() - due_date)) = -2").
			Find(&responses).Error
		if err != nil {
			log.Printf("SendNotifications cron executed. error : %v", err)
			return
		}
		for _, res := range responses {
			log.Printf("%+v\n", res)

			type LoanPaymentStatus struct {
				ID            uint64 `json:"id"`
				PaymentStatus string `json:"payment_status"`
			}
			var loanStatus LoanPaymentStatus
			err = DB.Table("loans").
				Where("? IN (SELECT UNNEST(loans.installment_details) )", res.ID).
				Find(&loanStatus).Error
			if err != nil {
				log.Printf("installment dont have valid parent (loan) id")
				continue
			}
			log.Printf("SendNotifications cron executed. loanStatus : %+v ", loanStatus)
		}

		log.Printf("SendNotifications cron executed. response : %+v ", responses)
	}
}
