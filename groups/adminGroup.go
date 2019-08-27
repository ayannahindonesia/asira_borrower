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

	// config info
	g.GET("/info", handlers.AsiraAppInfo)

	// Images
	g.GET("/image/:image_id", admin_handlers.GetImageB64String)

	//Borrowers
	g.GET("/borrower", admin_handlers.BorrowerGetAll)
	g.GET("/borrower/:borrower_id", admin_handlers.BorrowerGetDetails)

	//Loans
	g.GET("/loan", admin_handlers.LoanGetAll)
	g.GET("/loan/:loan_id", admin_handlers.LoanGetDetails)
}
