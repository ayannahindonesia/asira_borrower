package groups

import (
	"asira/handlers"
	"asira/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g)
	g.POST("/register_borrower", handlers.RegisterBorrower)
}
