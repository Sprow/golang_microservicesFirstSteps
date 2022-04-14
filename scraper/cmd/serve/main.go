package main

import (
	"context"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"scrapy/internal/manager"
	"scrapy/internal/mongodb"
	"scrapy/internal/scrap"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://sprow:12345@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//q, err := ch.QueueDeclare(
	//	"hello", // name
	//	false,   // durable
	//	false,   // delete when unused
	//	false,   // exclusive
	//	false,   // no-wait
	//	nil,     // arguments
	//)
	//failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		"urls", // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Println(err)
	}
	collectionChat := client.Database("siteData").Collection("data")

	httpClient := &http.Client{}
	s := scrap.NewSiteScraper(httpClient)
	db := mongodb.NewMongoDB(collectionChat)
	m := manager.NewSiteScraperManager(s, db)

	var forever chan struct{}

	go func() {
		for d := range msgs {
			err = m.GetSiteData(ctx, string(d.Body))
			if err != nil {
				log.Println(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
