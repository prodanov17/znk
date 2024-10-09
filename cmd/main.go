package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/prodanov17/znk/cmd/api"
)

func main() {

	server := api.NewAPIServer(fmt.Sprintf(":%s", "8000"), nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
