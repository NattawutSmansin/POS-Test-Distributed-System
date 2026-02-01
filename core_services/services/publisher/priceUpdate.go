package publisher

import (
	"encoding/json"
	"fmt"
	"log"

	"cor-service/requests"
	rabbitmq "cor-service/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PriceUpdateExchange = "price_updates"
)

type PriceUpdatePublisher struct {
	conn     *rabbitmq.Connection
	confirms chan amqp.Confirmation
}

func NewPriceUpdatePublisher(conn *rabbitmq.Connection) *PriceUpdatePublisher {
	return &PriceUpdatePublisher{conn: conn}
}

func (p *PriceUpdatePublisher) Setup() error {
	ch := p.conn.Channel

	err := ch.ExchangeDeclare(
		PriceUpdateExchange, // name
		"fanout",            // type
		true,                // durable
		false,               // auto-delete
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	err = ch.Confirm(false)
	if err != nil {
		return fmt.Errorf("failed to set confirm mode: %w", err)
	}

	p.confirms = ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	log.Println("setup complete")
	return nil
}

func (p *PriceUpdatePublisher) Publish(update requests.ProductPrice) error {
	ch := p.conn.Channel

	body, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal price update: %w", err)
	}

	err = ch.Publish(
		PriceUpdateExchange, // exchange
		"",                  // routing key: fanout ไม่ใช้
		true,                // mandatory
		false,               // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // survive restart
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish price update: %w", err)
	}

	confirm, ok := <-p.confirms
	if !ok {
		return fmt.Errorf("confirm channel closed before confirmation received")
	}
	if !confirm.Ack {
		return fmt.Errorf("broker nack message with delivery tag %d", confirm.DeliveryTag)
	}

	log.Printf(
		"published price update: product_id=%d, new_price=%.2f",
		update.ID,
		update.Price,
	)

	return nil
}
