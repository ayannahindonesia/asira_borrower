package groups

import (
	"asira_borrower/handlers"
	"asira_borrower/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "client")
	g.GET("/check_unique", handlers.CheckData)
	g.POST("/register_borrower", handlers.RegisterBorrower)
	g.POST("/borrower_login", handlers.BorrowerLogin)
	g.GET("/imagefile/:file_id", handlers.ClientImageFile)

	g.POST("/reset_password", handlers.ClientResetPassword)
	g.POST("/change_password", handlers.ChangePassword)

	//banks
	g.GET("/banks", handlers.ClientBanks)
	g.GET("/banks/:bank_id", handlers.ClientBankbyID)
}
