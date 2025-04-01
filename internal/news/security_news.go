package news

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// RSS represents a simple RSS feed.
type RSS struct {
	Channel struct {
		Title string `xml:"title"`
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

// FetchSecurityNews retrieves and processes security news from an RSS feed.
func FetchSecurityNews() {
	url := "https://www.securityweek.com/rss" // Example RSS feed
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Failed to fetch security news: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read RSS feed: %v", err)
		return
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		log.Printf("Failed to parse RSS feed: %v", err)
		return
	}

	log.Printf("Security News Feed: %s", rss.Channel.Title)
	for _, item := range rss.Channel.Items {
		log.Printf("News: %s - %s", item.Title, item.Link)
	}
}
