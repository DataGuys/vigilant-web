package crawler

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"golang.org/x/net/proxy"
)

func CrawlOnionSite(onionURL, torProxy string) (map[string]string, error) {
	results := make(map[string]string)

	// Provide default or fallback
	if torProxy == "" {
		torProxy = "127.0.0.1:9050"
	}

	dialer, err := proxy.SOCKS5("tcp", torProxy, nil, &net.Dialer{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("could not init SOCKS5 dialer: %w", err)
	}
	httpTransport := &http.Transport{
		DialContext: dialer.DialContext,
		// Optionally add TLS settings if needed
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
	)
	c.WithTransport(httpTransport)
	c.SetRequestTimeout(45 * time.Second)

	c.OnHTML("html", func(e *colly.HTMLElement) {
		results[e.Request.URL.String()] = e.Text
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("[CrawlOnion] Error on %s: %v", r.Request.URL, err)
	})

	if err := c.Visit(onionURL); err != nil {
		return results, fmt.Errorf("visit error: %w", err)
	}
	c.Wait()

	return results, nil
}
