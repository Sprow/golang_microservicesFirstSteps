package main

import (
	"ContentScraper/cmd/serve/handler"
	"ContentScraper/internal/manager"
	"ContentScraper/internal/mongodb"
	"ContentScraper/internal/scraper"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
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

//func connectRabbitMQ(conn *amqp.Connection) {
//	for i := 0; i < 20; i++ {
//		time.Sleep(2 * time.Second)
//		var err error
//		conn, err = amqp.Dial("amqp://sprow:12345@rabbitmq:5672/")
//		if err == nil {
//			break
//		}
//		failOnError(err, "Failed to connect to RabbitMQ")
//	}
//}

func main() {
	//conn, err := amqp.Dial("amqp://sprow:12345@localhost:5672/")  // local use
	conn, err := amqp.Dial("amqp://sprow:12345@rabbitmq:5672/") // run in docker
	failOnError(err, "Failed to connect to RabbitMQ")
	//var conn *amqp.Connection
	//connectRabbitMQ(conn)
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

	urls, err := ch.Consume(
		"urls", // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	err = ch.ExchangeDeclare( // Create exchange to send scraped sites _id to content_parser
		"parser_direct", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx := context.Background()
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")) // use for localhost
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017")) // for docker run
	if err != nil {
		log.Println(err)
	}
	collectionSiteContent := client.Database("siteContent").Collection("content")

	httpClient := &http.Client{}
	s := scraper.NewContentScraper(httpClient)
	db := mongodb.NewMongoDB(collectionSiteContent)
	m := manager.NewContentScraperManager(ch, s, db)

	h := handler.NewHandler(m)
	router := chi.NewRouter()
	h.Register(router)

	go func() {
		for d := range urls {
			err = m.GetSiteContent(ctx, string(d.Body))
			if err != nil {
				log.Println(err)
				//d.Ack(true) // task не выполнен, возвращаем его в rabbitMQ
				//continue
			}
			d.Ack(false)            // подтверждаем что task выполнен
			time.Sleep(time.Second) // для проверки как работают несколько скраперов
		}
	}()

	err = http.ListenAndServe(":8082", router)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Waiting for messages.")
}
