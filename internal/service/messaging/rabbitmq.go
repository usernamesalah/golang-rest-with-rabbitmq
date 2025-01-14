package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Messagging interface {
	Publish(ctx context.Context, queueName string, message interface{}) error
	Consume(ctx context.Context, queueName string, handler func(string)) error
	CreateQueue(queueName string) error
	DeleteQueue(queueName string) error
	StartQueue(ctx context.Context, queueName string, handler func(string)) error
}

type RabbitMQ struct {
	channel *amqp091.Channel
}

// NewRabbitMQ initializes a new RabbitMQ instance
func NewRabbitMQ(mqConn *amqp091.Connection) Messagging {
	mqChannel, err := mqConn.Channel()
	if err != nil {
		logrus.Fatalf("failed to create RabbitMQ channel: %s", err)
	}

	return &RabbitMQ{
		channel: mqChannel,
	}
}

// Publish sends a message to the specified queue
func (mq *RabbitMQ) Publish(ctx context.Context, queueName string, message interface{}) error {
	_, err := mq.channel.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	bodyJson, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = mq.channel.PublishWithContext(
		ctx,
		"", // Exchange
		queueName, false, false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        bodyJson,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Consume sets up a consumer for the specified queue
func (mq *RabbitMQ) Consume(ctx context.Context, queueName string, handler func(string)) error {
	msgs, err := mq.channel.ConsumeWithContext(
		ctx, queueName, "", true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			handler(string(msg.Body))
		}
	}()

	return nil
}

// CreateQueue explicitly creates a RabbitMQ queue
func (mq *RabbitMQ) CreateQueue(queueName string) error {
	_, err := mq.channel.QueueDeclare(
		queueName, true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create queue '%s': %w", queueName, err)
	}
	return nil
}

// DeleteQueue deletes a RabbitMQ queue with the given name
func (mq *RabbitMQ) DeleteQueue(queueName string) error {
	_, err := mq.channel.QueueDelete(
		queueName, // Name of the queue
		false,     // IfUnused: only delete if the queue has no consumers
		false,     // IfEmpty: only delete if the queue is empty
		false,     // No-wait: does not wait for server confirmation
	)
	if err != nil {
		return fmt.Errorf("failed to delete queue '%s': %w", queueName, err)
	}
	return nil
}

// StartQueue creates a queue (if not exists) and sets up a consumer
func (mq *RabbitMQ) StartQueue(ctx context.Context, queueName string, handler func(string)) error {
	// Create the queue if it does not exist
	if err := mq.CreateQueue(queueName); err != nil {
		return fmt.Errorf("failed to start queue: %w", err)
	}

	// Start the consumer
	if err := mq.Consume(ctx, queueName, handler); err != nil {
		return fmt.Errorf("failed to start consumer for queue '%s': %w", queueName, err)
	}

	return nil
}
