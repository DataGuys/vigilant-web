package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/yourusername/vigilant-web/config"
	"github.com/yourusername/vigilant-web/internal/db"
	"github.com/yourusername/vigilant-web/internal/server"
)

func main() {
	// Setup global logger
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.InfoLevel)

	cfg := config.LoadConfig()

	database, err := db.InitDB(cfg.DBPath)
	if err != nil {
		logrus.Fatalf("DB init error: %v", err)
	}
	defer database.Close()

	srv := server.NewServer(database, cfg.HTTPPort, cfg.TorProxy, cfg.JWTSecret)
	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
