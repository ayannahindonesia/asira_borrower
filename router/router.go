package router

import (
	"asira/groups"
	"asira/handlers"

	"github.com/labstack/echo"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	// e.GET("/test", handlers.Test)
	e.GET("/clientauth", handlers.ClientLogin)

	groups.AdminGroup(e)
	groups.ClientGroup(e)
	groups.BorrowerGroup(e)
	groups.UnverifiedBorrowerGroup(e)

	return e
}
