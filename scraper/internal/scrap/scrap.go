package scrap

import (
	"io"
	"log"
	"net/http"
)

type SiteScraper struct {
	client *http.Client
}

func NewSiteScraper(c *http.Client) *SiteScraper {
	return &SiteScraper{
		client: c,
	}
}

func (s SiteScraper) GetData(url string) ([]byte, error) {
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
