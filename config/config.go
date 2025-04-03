package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath        string
	TorProxy      string
	HTTPPort      int
	JWTSecret     string
}

func LoadConfig() *Config {
	// load .env if it exists
	_ = godotenv.Load()

	dbPath := getEnv("DB_PATH", "vigilant.db")
	torProxy := getEnv("TOR_PROXY", "127.0.0.1:9050")
	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
	jwtSecret := getEnv("JWT_SECRET", "change-this-secret")

	return &Config{
		DBPath:    dbPath,
		TorProxy:  torProxy,
		HTTPPort:  httpPort,
		JWTSecret: jwtSecret,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
