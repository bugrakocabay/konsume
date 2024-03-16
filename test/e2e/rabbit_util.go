package e2e

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

// connectToRabbitMQ establishes a connection to RabbitMQ and returns the connection and channel
func connectToRabbitMQ(connectionString string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

// pushMessageToQueue publishes a message to the specified queue in RabbitMQ
func pushMessageToQueue(ch *amqp.Channel, queueName string, body []byte) error {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}
