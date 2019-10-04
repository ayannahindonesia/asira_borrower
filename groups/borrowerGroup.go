package groups

import (
	"asira_borrower/handlers"
	"asira_borrower/middlewares"

	"github.com/labstack/echo"
)

func BorrowerGroup(e *echo.Echo) {
	g := e.Group("/borrower")
	middlewares.SetClientJWTmiddlewares(g, "borrower")

	// Profile endpoints
	g.GET("/profile", handlers.BorrowerProfile)
	g.PATCH("/profile", handlers.BorrowerProfileEdit)

	// Loan endpoints
	g.GET("/loan", handlers.BorrowerLoanGet)
	g.GET("/loan/:loan_id/details", handlers.BorrowerLoanGetDetails)
	g.POST("/loan", handlers.BorrowerLoanApply)
	g.GET("/loan/:loan_id/otp", handlers.BorrowerLoanOTPrequest)
	g.POST("/loan/:loan_id/verify", handlers.BorrowerLoanOTPverify)

	// Bank Endpoint
	g.GET("/bank_services", handlers.BorrowerBankService)
	g.GET("/bank_services/:service_id", handlers.BorrowerBankServiceDetails)
	// g.GET("/bank_products", handlers.BorrowerBankProduct)
	// g.GET("/bank_products/:product_id", handlers.BorrowerBankProductDetails)
}

func UnverifiedBorrowerGroup(e *echo.Echo) {
	g := e.Group("/unverified_borrower")
	middlewares.SetClientJWTmiddlewares(g, "unverified_borrower")

	// OTP
	g.POST("/otp_request", handlers.RequestOTPverifyAccount)
	g.POST("/otp_verify", handlers.VerifyAccountOTP)
}
