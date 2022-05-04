package scraper

import (
	"io"
	"log"
	"net/http"
)

type ContentScraper struct {
	client *http.Client
}

func NewContentScraper(c *http.Client) *ContentScraper {
	return &ContentScraper{
		client: c,
	}
}

func (s ContentScraper) GetContent(url string) ([]byte, error) {
	var bodyBytes []byte
	resp, err := s.client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return bodyBytes, err
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return bodyBytes, err
		}
	}
	return bodyBytes, nil
}
