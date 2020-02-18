package groups

import (
	"asira_borrower/handlers"
	"asira_borrower/middlewares"

	"github.com/labstack/echo"
)

func AgentGroup(e *echo.Echo) {
	g := e.Group("/agent")
	middlewares.SetClientJWTmiddlewares(g, "agent")

	// agent's profile endpoints
	g.GET("/profile", handlers.AgentProfile)
	g.PATCH("/profile", handlers.AgentProfileEdit)

	// agent's profile endpoints
	g.POST("/register_borrower", handlers.AgentRegisterBorrower)

	// OTP after register
	g.POST("/otp_request/:borrower_id", handlers.AgentRequestOTP)
	g.POST("/otp_verify/:borrower_id", handlers.AgentVerifyOTP)

	//banks owned by current agent (jti)
	g.GET("/banks", handlers.AgentAllBank)

	// agent's bank Endpoint
	g.GET("/bank_services", handlers.AgentBankService)

	// agent's bank Endpoint
	g.GET("/bank_products", handlers.AgentBankProduct)

	//borrowers owned by current agent (jti) and bank_id
	g.GET("/borrowers", handlers.AgentAllBorrower)

	//borrowers owned by current agent (jti) and agent's borrower_id
	g.GET("/borrower/:borrower_id", handlers.AgentBorrowerProfile)
	g.PATCH("/borrower/:borrower_id", handlers.AgentBorrowerProfileEdit)

	//check borrower from agent is exist or not
	g.POST("/checks_borrower", handlers.AgentCheckBorrower)

	// Loan endpoints
	g.POST("/loan", handlers.AgentLoanApply)
	g.GET("/loan", handlers.AgentLoanGet)
	g.GET("/loan/:loan_id/details", handlers.AgentLoanGetDetails)
	g.GET("/loan/:loan_id/otp", handlers.AgentLoanOTPrequest)
	g.POST("/loan/:loan_id/verify", handlers.AgentLoanOTPverify)

}
