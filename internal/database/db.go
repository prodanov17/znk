package database

import (
	"database/sql"
	"fmt"
	"log"
)

func NewPGStorage(cfg *PGConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

type PGConfig struct {
	DBUser   string
	DBPasswd string
	DBHost   string
	DBPort   string
	DBName   string
	SSLMode  string
}

func (pg *PGConfig) FormatDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", pg.DBUser, pg.DBPasswd, pg.DBHost, pg.DBPort, pg.DBName, pg.SSLMode)
}
