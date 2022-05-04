package parser

import (
	"bytes"
	"content_parser/internal/mongodb"
	"context"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

type ContentParser struct {
	db *mongodb.MongoDB
}

func NewContentParser(db *mongodb.MongoDB) *ContentParser {
	return &ContentParser{
		db: db,
	}
}

func (p *ContentParser) ParseAndSave(content []byte) error {
	reader := bytes.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
		return err
	}
	weather := doc.Find(".today-temperature").Text()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//ctx := context.Background()
	err = p.db.Save(ctx, weather)

	return nil
}
