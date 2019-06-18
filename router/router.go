package router

import (
	"asira/groups"
	"asira/handlers"

	"github.com/labstack/echo"
)

func NewBorrower() *echo.Echo {
	e := echo.New()

	// e.GET("/test", handlers.Test)
	e.GET("/clientauth", handlers.ClientLogin)

	groups.ClientGroup(e)
	groups.BorrowerGroup(e)

	return e
}
