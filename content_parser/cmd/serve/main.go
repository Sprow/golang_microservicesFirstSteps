package main

import (
	"content_parser/internal/mongodb"
	"content_parser/internal/parser"
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"net/url"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type HtmlBody struct {
	Html []byte `json:"html"`
}

func getContentBody(id amqp.Delivery) ([]byte, error) {
	data := make(url.Values)
	data.Add("oid", string(id.Body))

	resp, err := http.PostForm("http://content_scraper:8082", data)
	if err != nil {
		log.Println("failed to get data", err)
		return []byte{}, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	var content HtmlBody
	err = d.Decode(&content)
	return content.Html, err
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
	conn, err := amqp.Dial("amqp://sprow:12345@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"oid", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.ExchangeDeclare(
		"parser_direct", // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.QueueBind(
		q.Name,          // queue name
		"oid",           // routing key
		"parser_direct", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	IDs, err := ch.Consume(
		"oid", // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	// mongoDB
	ctx := context.Background()
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")) // use for localhost
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017")) // for docker run
	if err != nil {
		log.Println(err)
	}
	collectionWeather := client.Database("Weather").Collection("kyiv")
	db := mongodb.NewMongoDB(collectionWeather)
	p := parser.NewContentParser(db)

	go func() {
		for id := range IDs {
			content, err := getContentBody(id)
			if err != nil {
				log.Println(err)
				return
			}
			err = p.ParseAndSave(content)
			if err != nil {
				log.Println("fail ParseAndSave in main")
				log.Println(err)
				return
			}
			id.Ack(false) // подтверждаем что task выполнен
		}
	}()

	log.Println("Waiting for content id.")
	select {}
}
