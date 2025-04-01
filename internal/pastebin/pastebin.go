package pastebin

import (
	"log"
	"net/http"
	"time"
	"io/ioutil"
	"regexp"
)

// MonitorPastebin fetches recent public pastes and looks for .onion links.
func MonitorPastebin() {
	url := "https://pastebin.com/archive" // Using the archive page as a proxy for public pastes.
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Failed to fetch Pastebin archive: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read Pastebin archive: %v", err)
		return
	}

	// Simple regex to match .onion URLs
	re := regexp.MustCompile(`https?://[a-z0-9]{16}\.onion`)
	matches := re.FindAllString(string(body), -1)
	log.Printf("Found %d .onion links in Pastebin archive.", len(matches))
	for _, link := range matches {
		log.Printf("Pastebin link: %s", link)
	}
}
