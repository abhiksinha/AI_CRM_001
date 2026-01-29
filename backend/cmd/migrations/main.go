package main

import (
	"CRM/packages/configloader"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]

	// --- Load Configuration ---
	cfg, err := configloader.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// --- Database Connection ---
	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// --- Run Migrations ---
	// The directory where your .go migration files are located.
	migrationsDir := "internal/migrations"

	fmt.Printf("Running goose command '%s' on directory '%s'\n", command, migrationsDir)

	if err := goose.Run(command, db, migrationsDir, args[1:]...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}

	fmt.Println("Goose command finished successfully.")
}
