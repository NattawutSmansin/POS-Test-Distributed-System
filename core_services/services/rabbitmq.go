package services

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Connection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Config  Config
}

func NewConnectionRabbitMQ() (*Connection, error) {
	cfg := LoadRabbitMQConfig()

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	var conn *amqp.Connection
	var err error

	// retry loop
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("connection failed (attempt %d/5): %s", i+1, err)
		time.Sleep(time.Duration(i+1) * 2 * time.Second) // exponential backoff
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after 5 attempts: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	log.Println("connected successfully")

	return &Connection{
		Conn:    conn,
		Channel: ch,
		Config:  cfg,
	}, nil
}

func (c *Connection) Close() error {
	if err := c.Channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	log.Println("connection closed")
	return nil
}

func (c *Connection) IsClosed() bool {
	return c.Conn.IsClosed()
}

func LoadRabbitMQConfig() Config {
	port, _ := strconv.Atoi(os.Getenv("RABBITMQ_PORT"))

	return Config{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     port,
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}
}
