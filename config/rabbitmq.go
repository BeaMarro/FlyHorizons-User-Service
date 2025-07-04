package config

import "github.com/rabbitmq/amqp091-go"

var RabbitMQClient *RabbitMQ

// Struct to store RabbitMQ connection and channel
type RabbitMQ struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}
