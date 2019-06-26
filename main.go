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

	flags.Usage = usage
	flags.Parse(os.Args[1:])
	args := flags.Args()

	migrationDir := "migration" // migration directory

	switch args[0] {
	default:
		flags.Usage()
		break
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
		break
	case "seed":
		migration.Seed()
		os.Exit(0)
		break
	case "truncate":
		err := migration.Truncate(args[1:])
		if err != nil {
			log.Fatalf("%v", err)
			flags.Usage()
		}
		os.Exit(0)
		break
	case "create":
		if err := goose.Run("create", nil, migrationDir, args[1:]...); err != nil {
			log.Fatalf("goose create: %v", err)
			flags.Usage()
		}
		return
	case "migrate": // command example : [app name] migrate up
		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatalf("goose set dialect : %v", err)
			flags.Usage()
		}

		dbconf := asira.App.Config.GetStringMap(fmt.Sprintf("%s.database", asira.App.ENV))
		connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s", dbconf["host"].(string), dbconf["username"].(string), dbconf["table"].(string), dbconf["sslmode"].(string), dbconf["password"].(string))

		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatalf("-connectionString=%q: %v\n", connectionString, err)
			flags.Usage()
		}

		if err := goose.Run(args[1], db, migrationDir, args[2:]...); err != nil {
			log.Fatalf("goose run: %v", err)
			flags.Usage()
		}
		break
	}
}

func usage() {
	usagestring := `
to run the app :
	[app_name] run [app_mode]
	example : asira run borrower

to update db :
	[app_name] migrate [goose_command]
	example : asira migrate up
	goose command lists:
		up                   Migrate the DB to the most recent version available
		up-by-one            Migrate the DB up by 1
		up-to VERSION        Migrate the DB to a specific VERSION
		down                 Roll back the version by 1
		down-to VERSION      Roll back to a specific VERSION
		redo                 Re-run the latest migration
		reset                Roll back all migrations
		status               Dump the migration status for the current DB
		version              Print the current version of the database
		create NAME [sql|go] Creates new migration file with the current timestamp
		fix                  Apply sequential ordering to migrations

database seeding : (development environment only)
	[app_name] seed
	example : asira seed

database truncate : (development environment only)
	[app_name] truncate [table(s)]
	example : asira truncate borrowers | asira truncate borrowers loans | asira truncate all
	replace [table] with 'all' to truncate all tables
	`

	log.Print(usagestring)
}
