package manager

import (
	"ContentScraper/internal/mongodb"
	"ContentScraper/internal/scraper"
	"context"
	"log"
)

type SiteScraperManager struct {
	scraper *scraper.SiteScraper
	db      *mongodb.MongoDB
}

func NewSiteScraperManager(scraper *scraper.SiteScraper, db *mongodb.MongoDB) *SiteScraperManager {
	return &SiteScraperManager{
		scraper: scraper,
		db:      db,
	}
}

func (m *SiteScraperManager) GetSiteContent(ctx context.Context, url string) error {
	data, err := m.scraper.GetContent(url)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Get content from url %s", url)
	err = m.db.SaveContent(ctx, data)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Save content from url %s", url)
	return nil
}
