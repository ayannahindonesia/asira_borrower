package router

import (
	"asira_borrower/groups"
	"asira_borrower/handlers"
	"os"

	"github.com/labstack/echo"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	// e.GET("/test", handlers.Test)
	e.GET("/clientauth", handlers.ClientLogin)

	// files url
	gopath, _ := os.Getwd()
	e.Static("/", gopath+"/assets")

	groups.AdminGroup(e)
	groups.ClientGroup(e)
	groups.BorrowerGroup(e)
	groups.UnverifiedBorrowerGroup(e)

	return e
}
