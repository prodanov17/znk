package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/prodanov17/znk/cmd/api"
	"github.com/prodanov17/znk/internal/config"
	"github.com/prodanov17/znk/internal/database"
	"github.com/prodanov17/znk/internal/queue"
	"github.com/prodanov17/znk/pkg/logger"
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
		logger.Log.Fatal(err)
	}

	queueConn, err := queue.NewRabbitMQ(&queue.RabbitMQConfig{
		Host:     config.Env.RabbitMQHost,
		Port:     config.Env.RabbitMQPort,
		User:     config.Env.RabbitMQUser,
		Password: config.Env.RabbitMQPassword,
	})
	if err != nil {
		logger.Log.Fatal(err)
	}

	defer queueConn.Close()
	defer db.Close()

	initStorage(db)

	if err := server.Run(); err != nil {
		logger.Log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB: Successfully connected!")
}
