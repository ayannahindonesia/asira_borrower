package groups

import (
	"asira/handlers"
	"asira/middlewares"

	"github.com/labstack/echo"
)

func BorrowerGroup(e *echo.Echo) {
	g := e.Group("/borrower")
	middlewares.SetClientJWTmiddlewares(g, "borrower")
	g.GET("/profile", handlers.BorrowerProfile)
}
