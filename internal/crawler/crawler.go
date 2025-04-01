package crawler

import (
	"log"
	"net/http"

	"github.com/cretz/bine/tor"
	"github.com/dgraph-io/badger/v3"
	"github.com/gocolly/colly/v2"
)

// Crawler holds the Tor client and database reference.
type Crawler struct {
	torClient *tor.Tor
	db        *badger.DB
}

// NewCrawler creates a new Crawler instance.
func NewCrawler(t *tor.Tor, db *badger.DB) *Crawler {
	return &Crawler{torClient: t, db: db}
}

// Crawl starts crawling from the provided seed URL.
func (c *Crawler) Crawl(seedURL string) error {
	// Obtain a SOCKS5 dialer from Tor.
	dialer, err := c.torClient.Dialer(nil, nil)
	if err != nil {
		return err
	}

	// Set up the Colly collector with a custom transport.
	collector := colly.NewCollector()
	collector.WithTransport(&http.Transport{
		Dial: dialer.Dial,
	})

	// On every HTML response, store the page content in the database.
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		content := e.Text
		err := c.db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(url), []byte(content))
		})
		if err != nil {
			log.Printf("Failed to store %s: %v", url, err)
		} else {
			log.Printf("Stored %s successfully", url)
		}
	})

	// Start crawling.
	return collector.Visit(seedURL)
}
