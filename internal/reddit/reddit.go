package reddit

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type RedditResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// FetchRedditData fetches recent posts from a target subreddit.
func FetchRedditData() {
	subreddit := "darknet" // Replace with desired subreddit
	url := "https://www.reddit.com/r/" + subreddit + "/new.json?limit=5"
	client := http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "vigilant-onion-bot")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch Reddit data: %v", err)
		return
	}
	defer resp.Body.Close()

	var redditResp RedditResponse
	if err := json.NewDecoder(resp.Body).Decode(&redditResp); err != nil {
		log.Printf("Failed to parse Reddit JSON: %v", err)
		return
	}

	for _, child := range redditResp.Data.Children {
		log.Printf("Reddit Post: %s - %s", child.Data.Title, child.Data.URL)
	}
}
