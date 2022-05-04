package main

import (
	"ContentTask/cmd/serve/handler"
	"ContentTask/internal/content_task"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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
	//conn, err := amqp.Dial("amqp://sprow:12345@localhost:5672/") // local use
	conn, err := amqp.Dial("amqp://sprow:12345@rabbitmq:5672/") // docker container
	failOnError(err, "Failed to connect to rabbitmq")
	//var conn *amqp.Connection
	//connectRabbitMQ(conn)
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"tasks_direct", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	router := chi.NewRouter()
	p := content_task.NewContentTask(ch)
	h := handler.NewHandler(p)
	h.Register(router)

	err = http.ListenAndServe(":8081", router)
	if err != nil {
		fmt.Println(err)
	}
}
