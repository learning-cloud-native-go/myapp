package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"myapp/config"
)

const (
	dialect     = "pgx"
	fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
)

var (
	flags = flag.NewFlagSet("migration", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	c := config.NewDB()
	dbString := fmt.Sprintf(fmtDBString, c.Host, c.Username, c.Password, c.DBName, c.Port)

	db, err := goose.OpenDBWithDriver(dialect, dbString)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	if err := goose.RunContext(context.Background(), command, db, *dir, args[1:]...); err != nil {
		log.Fatalf("migration %v: %v", command, err)
	}
}
