package groups

import (
	"asira/handlers"
	"asira/middlewares"

	"github.com/labstack/echo"
)

func BorrowerGroup(e *echo.Echo) {
	g := e.Group("/borrower")
	middlewares.SetClientJWTmiddlewares(g, "borrower")

	// Profile endpoints
	g.GET("/profile", handlers.BorrowerProfile)

	// Loan endpoints
	g.GET("/loan", handlers.BorrowerLoanGet)
	g.POST("/loan", handlers.BorrowerLoanApply)
}

func UnverifiedBorrowerGroup(e *echo.Echo) {
	g := e.Group("/unverified_borrower")
	middlewares.SetClientJWTmiddlewares(g, "unverified_borrower")

	// OTP
	g.POST("/otp_request", handlers.RequestOTPverifyAccount)
	g.POST("/otp_verify", handlers.VerifyAccountOTP)
}
