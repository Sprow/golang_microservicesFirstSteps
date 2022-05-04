package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
	"time"
)

type MongoDB struct {
	collection *mongo.Collection
}

func NewMongoDB(collection *mongo.Collection) *MongoDB {
	return &MongoDB{
		collection: collection,
	}
}

type data struct {
	OID     primitive.ObjectID `bson:"_id"`
	Weather string             `bson:"weather"`
	Time    time.Time          `bson:"time"`
}

func (db *MongoDB) Save(ctx context.Context, weather string) error {
	newID := primitive.NewObjectID()
	weather = strings.Split(weather, " ")[0] // delete spaces and newline \n
	var d = data{
		OID:     newID,
		Time:    time.Now(),
		Weather: weather,
	}
	_, err := db.collection.InsertOne(ctx, d)
	log.Printf("Save %s", weather)
	return err
}
