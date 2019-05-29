package main

import (
	"kayacredit/handlers"
	"os"

	"github.com/labstack/echo"
)

func main() {
	e := NewBorrower()
	e.Logger.Fatal(e.Start(":8000"))
	os.Exit(0)
}

func NewBorrower() *echo.Echo {
	e := echo.New()

	e.GET("/test", handlers.Test)

	return e
}
