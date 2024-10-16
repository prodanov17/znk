package queue

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type RabbitMQ struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewRabbitMQ(cfg *RabbitMQConfig) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (cfg *RabbitMQConfig) FormatDSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
}

func (r *RabbitMQ) Close() error {
	if err := r.ch.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}
