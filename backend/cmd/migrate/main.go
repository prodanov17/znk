package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/prodanov17/znk/internal/config"
	"github.com/prodanov17/znk/pkg/logger"
)

func main() {
	// Set up PostgreSQL connection string (DSN)
	dsn := "postgres://" + config.Env.DBUser + ":" + config.Env.DBPassword +
		"@" + config.Env.DBAddress + "/" + config.Env.DBName + "?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Ensure the database connection works
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	// Set up the PostgreSQL migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations", // Path to your migration files
		"znk",                           // Database name (PostgreSQL)
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get the current migration version
	v, d, _ := m.Version()
	log.Printf("Current version: %v, dirty: %v", v, d)

	// Handle migration commands (up/down)
	logger.Log.Info(os.Args)
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
