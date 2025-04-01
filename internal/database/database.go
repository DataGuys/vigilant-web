package database

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

// InitDB opens the BadgerDB with an encryption key.
func InitDB(dbPath string, encryptionKey []byte) (*badger.DB, error) {
	opts := badger.DefaultOptions(dbPath).
		WithEncryptionKey(encryptionKey).
		WithLogger(nil) // Disable Badger logging for production

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CloseDB safely closes the database.
func CloseDB(db *badger.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}
