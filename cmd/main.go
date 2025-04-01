package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"vigilant-onion/internal/crawler"
	"vigilant-onion/internal/database"
	"vigilant-onion/internal/darkweb"
	"vigilant-onion/internal/tor"
)

func main() {
	// Generate or load your 32-byte encryption key.
	// For production, use a secure key management process.
	encryptionKey := make([]byte, 32)
	_, err := rand.Read(encryptionKey)
	if err != nil {
		log.Fatalf("Failed to generate encryption key: %v", err)
	}
	log.Printf("Encryption Key: %s", hex.EncodeToString(encryptionKey))

	// Initialize BadgerDB (encrypted at rest)
	db, err := database.InitDB("./badger_db", encryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	// Start Tor client using Bine
	torClient, err := tor.StartTor()
	if err != nil {
		log.Fatalf("Failed to start Tor: %v", err)
	}
	defer torClient.Close()

	// Initialize the crawler with the Tor client and database
	c := crawler.NewCrawler(torClient, db)

	// Set up HTTP routes
	mux := http.NewServeMux()

	// /monitor endpoint to trigger darkweb monitoring/discovery
	mux.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {
		// Start darkweb discovery using the crawler and darkweb module
		seedURL := "http://example.onion"
		err := darkweb.Discover(seedURL, c, db)
		if err != nil {
			http.Error(w, "Monitoring failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Darkweb monitoring executed successfully."))
	})

	// Configure HTTPS server (data in transit encrypted)
	server := &http.Server{
		Addr:         ":443",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	log.Println("Starting HTTPS server on :443")
	// Replace with paths to your certificate and private key
	if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
