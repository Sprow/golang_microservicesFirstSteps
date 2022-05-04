package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	collection *mongo.Collection
}

func NewMongoDB(collection *mongo.Collection) *MongoDB {
	return &MongoDB{
		collection: collection,
	}
}

const (
	StatusScrapDone   = true
	StatusScrapFailed = false
)

type SiteContent struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	URL           string             `bson:"url" json:"URL"`
	Content       []byte             `bson:"content" json:"content"`
	ScrapedStatus bool               `bson:"scrapedStatus" json:"scrapedStatus"`
}

func (db *MongoDB) SaveContent(ctx context.Context, url string, content []byte, status bool) (primitive.ObjectID, error) {
	newID := primitive.NewObjectID()
	d := SiteContent{
		ID:            newID,
		URL:           url,
		Content:       content,
		ScrapedStatus: status,
	}
	_, err := db.collection.InsertOne(ctx, d)
	return newID, err
}

type HtmlBody struct {
	Html []byte `json:"html"`
}

func (db *MongoDB) GetContentByOID(ctx context.Context, oid primitive.ObjectID) (HtmlBody, error) {
	var content SiteContent
	filter := bson.D{{"_id", oid}}
	err := db.collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		return HtmlBody{content.Content}, err
	}
	return HtmlBody{content.Content}, nil
}
