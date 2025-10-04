package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// go run main.go --direction=up --dsn="postgres://user:pass@localhost:5432/mydb?sslmode=disable"
func main() {
	fmt.Println("🚀 Migrater: starting...")

	direction := flag.String("direction", "up", "Migration direction: up or down")
	step := flag.Int("step", 0, "Number of steps for partial migration (0 = all)")
	dsn := flag.String("dsn", "postgres://account_user:account_pass@localhost:5432/account_db?sslmode=disable", "PostgreSQL DSN")
	migrationsPath := flag.String("path", "migrations", "Path to migrations folder")
	flag.Parse()

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Cannot connect to DB: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migrationsPath),
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("❌ Failed to initialize migrate instance: %v", err)
	}

	switch *direction {
	case "up":
		if *step > 0 {
			fmt.Printf("⬆️  Migrating up %d step(s)...\n", *step)
			err = m.Steps(*step)
		} else {
			fmt.Println("⬆️  Migrating all pending migrations...")
			err = m.Up()
		}

	case "down":
		if *step > 0 {
			fmt.Printf("⬇️  Migrating down %d step(s)...\n", *step)
			err = m.Steps(-*step)
		} else {
			fmt.Println("⬇️  Rolling back last migration...")
			err = m.Steps(-1)
		}

	default:
		fmt.Println("⚠️  Invalid direction. Use --direction=up or --direction=down")
		os.Exit(1)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Migration error: %v", err)
	}

	fmt.Println("✅ Migration completed successfully.")
}
