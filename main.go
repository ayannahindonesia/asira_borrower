package main

import (
	"asira/asira"
	"asira/migration"
	"asira/router"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	defer asira.App.Close()

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
	case "goose": // command example : [app name] goose up
		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatalf("goose set dialect : %v", err)
		}

		dbconf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.database", asira.App.ENV))
		connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s", dbconf["host"].(string), dbconf["username"].(string), dbconf["table"].(string), dbconf["sslmode"].(string), dbconf["password"].(string))

		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatalf("-connectionString=%q: %v\n", connectionString, err)
		}

		if err := goose.Run(args[1], db, migrationDir, args[2:]...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		break
	}
}
