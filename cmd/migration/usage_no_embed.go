//go:build !embed

package main

import "fmt"

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migration COMMAND
Examples:
  migration create initial_tables sql
`

	usageCommands = `
Commands:
  create NAME [sql|go] Creates new migration file with the current timestamp
  fix                  Apply sequential ordering to migrations
  up                   Migrate the DB to the most recent version available
  up-by-one            Migrate the DB up by 1
  up-to VERSION        Migrate the DB to a specific VERSION
  down                 Roll back the version by 1
  down-to VERSION      Roll back to a specific VERSION
  redo                 Re-run the latest migration
  reset                Roll back all migrations
  status               Dump the migration status for the current DB
  version              Print the current version of the database`
)
