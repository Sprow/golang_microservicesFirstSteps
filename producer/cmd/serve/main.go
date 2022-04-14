package main

import (
	"Producer/cmd/serve/handler"
	"Producer/internal/producer"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://sprow:12345@localhost:5672/")
	failOnError(err, "Failed to connect to rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare( //declare a queue for us to send to
		"urls", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	router := chi.NewRouter()
	p := producer.NewProducer(ch)
	h := handler.NewHandler(p)
	h.Register(router)

	err = http.ListenAndServe(":8085", router)
	if err != nil {
		fmt.Println(err)
	}
}
