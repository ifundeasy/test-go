package queue

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ struct to hold the connection
type RabbitMQ struct {
	Conn *amqp.Connection
}

// NewRabbitMQ initializes a new RabbitMQ connection
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{Conn: conn}, nil
}

// PublishProductCreated publishes a message when a product is created
func (r *RabbitMQ) PublishProductCreated(product interface{}) error {
	return r.publish("product.created", product)
}

// PublishProductUpdated publishes a message when a product is updated
func (r *RabbitMQ) PublishProductUpdated(product interface{}) error {
	return r.publish("product.updated", product)
}

// PublishProductDeleted publishes a message when a product is deleted
func (r *RabbitMQ) PublishProductDeleted(productID string) error {
	return r.publish("product.deleted", productID)
}

// publish is a helper method to publish messages to RabbitMQ
func (r *RabbitMQ) publish(routingKey string, message interface{}) error {
	ch, err := r.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}

	log.Printf("Message published: %s", routingKey)
	return nil
}

// GetRabbitMQ initializes a new RabbitMQ connection and returns a RabbitMQ instance
func GetRmqInstance(rabbitMQURL string) *RabbitMQ {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	return &RabbitMQ{Conn: conn}
}

// Close closes the RabbitMQ connection
func (r *RabbitMQ) Close() error {
	return r.Conn.Close()
}
