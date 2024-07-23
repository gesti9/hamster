package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchTimestamp() (int64, error) {
	res, err := http.Get("https://www.unixtimestamp.com/")
	if err != nil {
		return 0, fmt.Errorf("error fetching the page: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error loading HTTP response body: %w", err)
	}

	var timestamp int64
	doc.Find("div#main-segment div.value.epoch").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		_, err := fmt.Sscanf(strings.TrimSpace(text), "%d", &timestamp)
		if err != nil {
			log.Printf("error parsing timestamp: %v", err)
		}
	})

	if timestamp == 0 {
		return 0, fmt.Errorf("could not find timestamp on the page")
	}

	return timestamp, nil
}

func main() {
	timestamp, err := fetchTimestamp()
	if err != nil {
		log.Fatalf("Failed to fetch timestamp: %v", err)
	}

	url := "https://api.hamsterkombatgame.io/clicker/tap"
	requestBody := map[string]interface{}{
		"count":         1,
		"availableTaps": 150,
		"timestamp":     timestamp,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 17217114320315CXUPGAqUR9iCszVL2i5SD4pYe83ipCg3xTVGgvB1Msi9nF2C7RJoYHgwYVuoxOD452639799")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request successful!")
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}
}
