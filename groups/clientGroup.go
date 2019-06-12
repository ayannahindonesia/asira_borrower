package groups

import (
	handlers "kayacredit/handlers/borrower"
	"kayacredit/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g)
	g.POST("/register_borrower", handlers.RegisterBorrower)
}
