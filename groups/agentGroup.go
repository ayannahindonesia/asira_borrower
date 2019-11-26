package groups

import (
	"asira_borrower/handlers"
	"asira_borrower/middlewares"

	"github.com/labstack/echo"
)

func AgentGroup(e *echo.Echo) {
	g := e.Group("/agent")
	middlewares.SetClientJWTmiddlewares(g, "agent")

	// agent's profile endpoints
	g.GET("/profile", handlers.AgentProfile)

	// agent's profile endpoints
	g.POST("/register_borrower", handlers.AgentRegisterBorrower)

	// agent's bank Endpoint
	g.GET("/bank_services", handlers.AgentBankService)
	g.GET("/bank_services/:service_id", handlers.AgentBankServiceDetails)
}
