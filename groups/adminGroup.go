package groups

import (
	"asira_borrower/handlers"
	"asira_borrower/middlewares"

	"github.com/labstack/echo"
)

func AdminGroup(e *echo.Echo) {
	g := e.Group("/admin")
	middlewares.SetClientJWTmiddlewares(g, "admin")

	// OTP
	g.GET("/info", handlers.AsiraAppInfo)

	//Loans
	g.GET("/loan", admin_handlers.LoanGetAll)
	g.GET("/loan/:loan_id", admin_handlers.LoanGetDetails)
}
