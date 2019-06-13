package main

import (
	"database/sql"
	"flag"
	"kayacredit/kc"
	"kayacredit/migration"
	"kayacredit/router"
	"log"
	"os"

	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	defer kc.App.Close()

	flags.Parse(os.Args[1:])
	args := flags.Args()

	migrationDir := "migration" // migration directory

	switch args[0] {
	case "run":
		switch args[1] {
		default:
			break
		case "borrower":
			e := router.NewBorrower()
			e.Logger.Fatal(e.Start(":8000"))
			os.Exit(0)
			break
		}
	case "seed":
		migration.Seed()
		os.Exit(0)
		break
	case "create":
		if err := goose.Run("create", nil, migrationDir, args[1:]...); err != nil {
			log.Fatalf("goose create: %v", err)
		}
		return
	case "postgres":
		if err := goose.SetDialect(args[0]); err != nil {
			log.Fatal(err)
		}
		db, err := sql.Open(args[0], args[1])
		if err != nil {
			log.Fatalf("-dbstring=%q: %v\n", args[1], err)
		}
		arguments := append([]string{}, args[3:]...)

		if err := goose.Run(args[3], db, migrationDir, arguments...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		break
	}
}
