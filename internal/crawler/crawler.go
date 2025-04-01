package crawler

import (
	"log"
	"net/http"

	"github.com/cretz/bine/tor"
	"github.com/dgraph-io/badger/v3"
	"github.com/gocolly/colly/v2"
)

// Crawler holds references to the Tor client and the database.
type Crawler struct {
	torClient *tor.Tor
	db        *badger.DB
}

// NewCrawler returns a new Crawler instance.
func NewCrawler(t *tor.Tor, db *badger.DB) *Crawler {
	return &Crawler{torClient: t, db: db}
}

// Crawl starts crawling from the given seed URL.
func (c *Crawler) Crawl(seedURL string) error {
	// Get a SOCKS5 dialer from the Tor client.
	dialer, err := c.torClient.Dialer(nil, nil)
	if err != nil {
		return err
	}

	// Initialize Colly collector with a custom HTTP transport (using Tor's dialer).
	collector := colly.NewCollector()
	collector.WithTransport(&http.Transport{
		Dial: dialer.Dial,
	})

	// On every HTML response, store the page content in the database.
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		content := e.Text

		// For demonstration, we store the URL and content as key/value.
		err := c.db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(url), []byte(content))
		})
		if err != nil {
			log.Printf("Failed to store %s: %v", url, err)
		} else {
			log.Printf("Stored %s successfully", url)
		}
	})

	// Start the crawl.
	return collector.Visit(seedURL)
}
