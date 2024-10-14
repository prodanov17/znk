package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host                   string
	Port                   string
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTSecret              string
	DeckPath               string
	JWTExpirationInSeconds int64
}

var Env = initConfig()

func initConfig() Config {
	godotenv.Load(".env")

	return Config{
		Host:                   getEnv("HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8000"),
		DBHost:                 getEnv("DB_HOST", "host.docker.internal"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPassword:             getEnv("DB_PASSWORD", ""),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "host.docker.internal"), getEnv("DB_PORT", "5432")),
		DBName:                 getEnv("DB_DATABASE", "znk"),
		JWTSecret:              getEnv("JWT_SECRET", ""),
		DeckPath:               getEnv("DECK_PATH", "assets/cards.json"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
