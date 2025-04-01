package main

import (
    "log"
    "vigilant-onion/internal/crawler"
    "vigilant-onion/internal/database"
    "vigilant-onion/internal/tor"
)

func main() {
    // Initialize the database
    db, err := database.InitDB("onionsites.db")
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    // Start the Tor client
    t, err := tor.StartTor()
    if err != nil {
        log.Fatalf("Failed to start Tor: %v", err)
    }
    defer t.Close()

    // Initialize the crawler with the Tor client and database
    c := crawler.NewCrawler(t, db)

    // Start crawling from a seed URL
    seedURL := "http://example.onion"
    if err := c.Crawl(seedURL); err != nil {
        log.Printf("Crawling failed: %v", err)
    }
}
