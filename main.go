package main

import (
	"kayacredit/handlers"
	"kayacredit/kc"
	"os"

	"github.com/labstack/echo"
)

var (
	kcApp kc.App
)

func main() {
	defer kcApp.Close()

	e := NewBorrower()
	e.Logger.Fatal(e.Start(":8000"))
	os.Exit(0)
}

func NewBorrower() *echo.Echo {
	e := echo.New()

	e.GET("/test", handlers.Test)

	return e
}
