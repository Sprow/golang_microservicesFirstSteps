package manager

import (
	"ContentScraper/internal/mongodb"
	"ContentScraper/internal/scraper"
	"context"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type ContentScraperManager struct {
	ch      *amqp.Channel
	scraper *scraper.ContentScraper
	DB      *mongodb.MongoDB
}

func NewContentScraperManager(ch *amqp.Channel, scraper *scraper.ContentScraper, db *mongodb.MongoDB) *ContentScraperManager {
	return &ContentScraperManager{
		ch:      ch,
		scraper: scraper,
		DB:      db,
	}
}

func (m *ContentScraperManager) GetSiteContent(ctx context.Context, url string) error {
	data, err := m.scraper.GetContent(url)
	var status bool
	if err != nil {
		log.Println(err)
		status = mongodb.StatusScrapFailed
		_, err = m.DB.SaveContent(ctx, url, []byte{}, status)
		return err
	}
	log.Printf("Get content from url %s", url)
	status = mongodb.StatusScrapDone
	oid, err := m.DB.SaveContent(ctx, url, data, status)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Save content from url %s", url)
	if status {
		err = m.publishScrapedSiteID(oid)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (m *ContentScraperManager) publishScrapedSiteID(oid primitive.ObjectID) error {
	err := m.ch.Publish(
		"parser_direct", // exchange
		"oid",           // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(oid.Hex()),
		})
	log.Printf("send _oid `%s` to exchange `parser_direct`, routing key `oid`", oid)
	return err
}
