package darkweb

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"vigilant-onion/internal/crawler"
)

// Discover starts the darkweb discovery process from the seed URL.
func Discover(seedURL string, c *crawler.Crawler, db *badger.DB) error {
	log.Printf("Starting darkweb discovery from %s", seedURL)
	if err := c.Crawl(seedURL); err != nil {
		return err
	}
	// Further link extraction and iterative crawling can be added here.
	log.Println("Darkweb discovery completed.")
	return nil
}
