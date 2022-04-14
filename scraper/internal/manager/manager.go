package manager

import (
	"context"
	"log"
	"scrapy/internal/mongodb"
	"scrapy/internal/scrap"
)

type SiteScraperManager struct {
	scraper *scrap.SiteScraper
	db      *mongodb.MongoDB
}

func NewSiteScraperManager(scraper *scrap.SiteScraper, db *mongodb.MongoDB) *SiteScraperManager {
	return &SiteScraperManager{
		scraper: scraper,
		db:      db,
	}
}

func (m *SiteScraperManager) GetSiteData(ctx context.Context, url string) error {
	data, err := m.scraper.GetData(url)
	if err != nil {
		log.Println(err)
		return err
	}
	err = m.db.SaveData(ctx, data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
