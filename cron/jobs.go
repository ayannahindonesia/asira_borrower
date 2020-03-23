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
			Where("extract('day' from date_trunc('day', NOW() - due_date)) = -3").
			Find(&responses).Error
		if err != nil {
			log.Printf("SendNotifications cron executed. error : %v", err)
			return
		}
		for _, res := range responses {
			log.Printf("%+v\n", res)

			// err = DB.Table("installments").
			// 	Where("extract('day' from date_trunc('day', NOW() - due_date)) = -3").
			// 	Find(&res).Error
			// if err != {
			// 	log.Printf("invalid id")
			// 	continue
			// }
		}
		// Where("disburse_date != ?", "0001-01-01 00:00:00+00").
		// Where("NOW() > disburse_date + make_interval(days => 2)").
		// Update("disburse_status", "confirmed").Error

		log.Printf("SendNotifications cron executed. response : %+v ", responses)
	}
}
