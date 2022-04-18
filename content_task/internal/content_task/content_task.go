package content_task

import (
	"github.com/streadway/amqp"
	"log"
)

type ContentTask struct {
	ch *amqp.Channel
}

func NewContentTask(ch *amqp.Channel) *ContentTask {
	return &ContentTask{
		ch: ch,
	}
}

func (p *ContentTask) PublishURL(url string) error {
	err := p.ch.Publish(
		"tasks_direct", // exchange
		"urls",         // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(url),
		})
	log.Printf("send url `%s` to exchange `tasks_direct`, routing key `urls`", url)
	return err
}
