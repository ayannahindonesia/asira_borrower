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

	//Borrowers
	g.GET("/borrower", handlers.BorrowerGetAll)
	g.GET("/borrower/:borrower_id/details", handlers.BorrowerGetDetails)
}
