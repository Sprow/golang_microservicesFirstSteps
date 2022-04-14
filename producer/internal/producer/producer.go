package producer

import "github.com/streadway/amqp"

type Producer struct {
	ch *amqp.Channel
}

func NewProducer(ch *amqp.Channel) *Producer {
	return &Producer{
		ch: ch,
	}
}

func (p *Producer) PublishURL(url string) error {
	err := p.ch.Publish(
		"",     // exchange
		"urls", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(url),
		})
	return err
}
