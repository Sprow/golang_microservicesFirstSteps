package main

import (
	"ContentScraper/internal/manager"
	"ContentScraper/internal/mongodb"
	"ContentScraper/internal/scraper"
	"context"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	//conn, err := amqp.Dial("amqp://sprow:12345@localhost:5672/")  // local use
	conn, err := amqp.Dial("amqp://sprow:12345@rabbitmq:5672/") // run in docker
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"urls", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,         // queue name
		"urls",         // routing key
		"tasks_direct", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

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
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")) // use for localhost
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017")) // for docker run
	if err != nil {
		log.Println(err)
	}
	collectionSiteContent := client.Database("siteContent").Collection("content")

	httpClient := &http.Client{}
	s := scraper.NewSiteScraper(httpClient)
	db := mongodb.NewMongoDB(collectionSiteContent)
	m := manager.NewSiteScraperManager(s, db)

	var forever chan struct{}

	go func() {
		for d := range msgs {
			err = m.GetSiteContent(ctx, string(d.Body))
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second) // для проверки как работают несколько скраперов
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
