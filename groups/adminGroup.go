package groups

import (
	"asira_borrower/admin_handlers"
	"asira_borrower/middlewares"
	"handlers"

	"github.com/labstack/echo"
)

func AdminGroup(e *echo.Echo) {
	g := e.Group("/admin")
	middlewares.SetClientJWTmiddlewares(g, "admin")

	// OTP
	g.GET("/info", handlers.AsiraAppInfo)

	//Create Client Config
	g.POST("/client_config", admin_handlers.CreateClientConfig)
}
