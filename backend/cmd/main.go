package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/prodanov17/znk/cmd/api"
	"github.com/prodanov17/znk/internal/config"
	"github.com/prodanov17/znk/internal/database"
)

func main() {

	server := api.NewAPIServer(fmt.Sprintf(":%s", "8000"), nil)
	db, err := database.NewPGStorage(&database.PGConfig{
		DBUser:   config.Env.DBUser,
		DBPasswd: config.Env.DBPassword,
		DBHost:   config.Env.DBHost,
		DBName:   config.Env.DBName,
		DBPort:   config.Env.DBPort,
		SSLMode:  "disable",
	})

	if err != nil {
		// log.Fatal(err)
	}

	initStorage(db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		// log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
