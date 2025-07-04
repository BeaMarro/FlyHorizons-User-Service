package health

import (
	"flyhorizons-userservice/config"
)

type RabbitMQCheck struct{}

func (c RabbitMQCheck) Name() string {
	return "rabbitmq"
}

func (c RabbitMQCheck) Pass() bool {
	conn := config.RabbitMQClient
	if conn == nil || conn.Connection == nil || conn.Channel == nil {
		return false
	}
	// Try to ping the connection
	return !conn.Connection.IsClosed()
}
