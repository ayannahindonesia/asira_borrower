package main

import (
	"kayacredit/kc"
	"kayacredit/router"
	"os"
)

func main() {
	defer kc.App.Close()

	e := router.NewBorrower()
	e.Logger.Fatal(e.Start(":8000"))
	os.Exit(0)
}
