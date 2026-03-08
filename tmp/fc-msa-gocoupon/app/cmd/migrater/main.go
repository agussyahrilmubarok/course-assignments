package main

import (
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbURL := flag.String("db", "", "Database connection URL")
	dir := flag.String("dir", "migrations", "Directory containing migration files")
	action := flag.String("action", "up", "Migration action: up | down | drop | force | version")
	forceVersion := flag.Int("force", 0, "Force set version (only works with --action=force)")
	flag.Parse()

	if *dbURL == "" {
		log.Fatal("missing required flag: --db")
	}

	m, err := migrate.New(
		"file://"+*dir,
		*dbURL,
	)

	if err != nil {
		log.Fatalf("failed to initialize migrate: %v", err)
	}

	switch *action {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "drop":
		err = m.Drop()
	case "force":
		err = m.Force(*forceVersion)
	case "version":
		v, dirty, _ := m.Version()
		log.Printf("version: %d, dirty: %v\n", v, dirty)
		return
	default:
		log.Fatalf("unknown action: %s", *action)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration error: %v", err)
	}

	log.Println("migration completed successfully")
}
