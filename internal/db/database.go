package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
)

var log = logrus.New()

func InitDB(dbPath string) (*sql.DB, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Infof("No existing DB at %s, will create new DB", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %w", err)
	}

	log.Infof("Connected to SQLite DB at %s", dbPath)

	// You can run your migrations here, e.g. with "goose" or "sql-migrate"
	// e.g. goose.Run("up", ...)

	return db, nil
}
