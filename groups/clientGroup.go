package groups

import (
	"asira/handlers"
	"asira/middlewares"

	"github.com/labstack/echo"
)

func ClientGroup(e *echo.Echo) {
	g := e.Group("/client")
	middlewares.SetClientJWTmiddlewares(g, "client")
	g.GET("/check_unique", handlers.CheckData)
	// e.GET("/ganteng", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "/users/:id")
	// })
	g.POST("/register_borrower", handlers.RegisterBorrower)
	g.POST("/borrower_login", handlers.BorrowerLogin)
}
