package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"cor-service/requests"
	rabbitmq "cor-service/services"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

const (
	OrderReportExchange   = "order_reports" // exchange เดียวกับที่ Branch publish
	OrderReportQueue      = "core_orders"   // queue สำหรับ Core รับ
	OrderReportRoutingKey = "order.report"  // routing key เดียวกับที่ Branch ใช้
)

type OrderReportConsumer struct {
	conn *rabbitmq.Connection
	repo *OrderRepository
	db   *gorm.DB
}

func NewOrderReportConsumer(conn *rabbitmq.Connection, database *gorm.DB) *OrderReportConsumer {
	repo := NewOrderRepository(database)
	return &OrderReportConsumer{conn: conn, db: database, repo: repo}
}

func (c *OrderReportConsumer) Setup() error {
	ch := c.conn.Channel

	err := ch.ExchangeDeclare(
		OrderReportExchange, // name
		"direct",            // type
		true,                // durable
		false,               // auto-delete
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = ch.QueueDeclare(
		OrderReportQueue, // name
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
		OrderReportQueue,      // queue name
		OrderReportRoutingKey, // routing key: "order.report"
		OrderReportExchange,   // exchange name
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

func (c *OrderReportConsumer) Start() error {
	ch := c.conn.Channel

	// prefetch 1
	err := ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := ch.Consume(
		OrderReportQueue, // queue
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

	log.Println("listening for order reports...")

	for msg := range msgs {
		if err := c.handleOrderReport(msg); err != nil {
			log.Printf("error processing message: %s", err)
			msg.Nack(false, true)
		} else {
			msg.Ack(false)
		}
	}

	return nil
}

func (c *OrderReportConsumer) handleOrderReport(msg amqp.Delivery) error {
	var report requests.OrderReport

	if err := json.Unmarshal(msg.Body, &report); err != nil {
		return fmt.Errorf("failed to unmarshal order report: %w", err)
	}

	log.Printf(
		"received order report: branch_id=%d, product_id=%d, total=%.2f",
		report.BranchID,
		report.ProductID,
		report.Total,
	)

	err := c.repo.CreateFromBranch(report)
	if err != nil {
		return fmt.Errorf("failed to save order to core_db: %w", err)
	}

	return nil
}
