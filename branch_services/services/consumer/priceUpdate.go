package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	rabbitmq "branch-service/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PriceUpdateExchange   = "price_updates" // exchange ที่ Core publish
	PriceUpdateQueue      = "branch_prices" // queue ของ branch นี้
	PriceUpdateRoutingKey = "price.update"  // routing key
)

type PriceUpdate struct {
	ID    int64   `json:"id"`
	Price float64 `json:"price"`
}

type PriceUpdateConsumer struct {
	conn *rabbitmq.Connection
}

func NewPriceUpdateConsumer(conn *rabbitmq.Connection) *PriceUpdateConsumer {
	return &PriceUpdateConsumer{conn: conn}
}

func (c *PriceUpdateConsumer) Setup() error {
	ch := c.conn.Channel

	err := ch.ExchangeDeclare(
		PriceUpdateExchange, // name
		"fanout",            // type: fanout = ส่งไปทุก queue ที่ bind อยู่
		true,                // durable: survive restart
		false,               // auto-delete
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = ch.QueueDeclare(
		PriceUpdateQueue, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		PriceUpdateQueue,    // queue name
		"",                  // routing key
		PriceUpdateExchange, // exchange name
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	log.Println("setup complete")
	return nil
}

func (c *PriceUpdateConsumer) Start() error {
	ch := c.conn.Channel

	// prefetch 1
	err := ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		PriceUpdateQueue, // queue
		"",               // consumer tag
		false,            // auto-ack: false = manual ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	log.Println("listening for price updates...")

	for msg := range msgs {
		if err := handlePriceUpdate(msg); err != nil {
			log.Printf("error processing message: %s", err)
			msg.Nack(false, true)
		} else {
			msg.Ack(false)
		}
	}

	return nil
}

func handlePriceUpdate(msg amqp.Delivery) error {
	var update PriceUpdate

	if err := json.Unmarshal(msg.Body, &update); err != nil {
		return fmt.Errorf("failed to unmarshal price update: %w", err)
	}

	log.Printf(
		"received price update: product_id=%d, new_price=%.2f",
		update.ID,
		update.Price,
	)

	return nil
}
