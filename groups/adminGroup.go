package groups

import (
	"asira_borrower/admin_handlers"
	"asira_borrower/handlers"
	"asira_borrower/middlewares"

	"github.com/labstack/echo"
)

func AdminGroup(e *echo.Echo) {
	g := e.Group("/admin")
	middlewares.SetClientJWTmiddlewares(g, "admin")

	// OTP
	g.GET("/info", handlers.AsiraAppInfo)

	//Borrowers
	g.GET("/borrower", admin_handlers.BorrowerGetAll)
	g.GET("/borrower/:borrower_id", admin_handlers.BorrowerGetDetails)

	//Loans
	g.GET("/loan", admin_handlers.LoanGetAll)
	g.GET("/loan/:loan_id", admin_handlers.LoanGetDetails)

	// Loan Purpose
	g.GET("/loan_purposes", admin_handlers.LoanPurposeList)
	g.POST("/loan_purposes", admin_handlers.LoanPurposeNew)
	g.GET("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposeDetail)
	g.PATCH("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposePatch)
	g.DELETE("/loan_purposes/:loan_purpose_id", admin_handlers.LoanPurposeDelete)
}
