package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emmiter struct {
	connection *amqp.Connection
}

func (e *Emmiter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emmiter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	log.Println("Pushing to channel")

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		return err
	}
	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emmiter, error) {
	emitter := Emmiter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emmiter{}, err
	}

	return emitter, nil

}
