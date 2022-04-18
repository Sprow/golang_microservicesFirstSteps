package mongodb

import (
	"context"
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

type siteData struct {
	ID   primitive.ObjectID `bson:"_id"`
	Data []byte             `bson:"data"`
}

func (db *MongoDB) SaveContent(ctx context.Context, data []byte) error {
	d := siteData{
		ID:   primitive.NewObjectID(),
		Data: data,
	}

	_, err := db.collection.InsertOne(ctx, d)
	return err
}
