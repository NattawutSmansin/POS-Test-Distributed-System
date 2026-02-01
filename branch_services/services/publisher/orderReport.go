package publisher

import (
	"encoding/json"
	"fmt"
	"log"

	"branch-service/requests"
	rabbitmq "branch-service/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	OrderReportExchange   = "order_reports"
	OrderReportRoutingKey = "order.report"
)

type OrderReportPublisher struct {
	conn     *rabbitmq.Connection
	confirms chan amqp.Confirmation
}

func NewOrderReportPublisher(conn *rabbitmq.Connection) *OrderReportPublisher {
	return &OrderReportPublisher{conn: conn}
}

func (p *OrderReportPublisher) Setup() error {
	ch := p.conn.Channel

	err := ch.ExchangeDeclare(
		OrderReportExchange, // name
		"direct",            // type: direct = route by routing key
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

	return nil
}

func (p *OrderReportPublisher) Publisher(report requests.OrderReport) error {
	ch := p.conn.Channel

	body, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal order report: %w", err)
	}

	err = ch.Publish(
		OrderReportExchange,   // exchange
		OrderReportRoutingKey, // routing key
		true,                  // mandatory: true = error ถ้า queue ไม่มี
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // survive restart
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish order report: %w", err)
	}

	confirm, ok := <-p.confirms
	if !ok {
		return fmt.Errorf("confirm channel closed before confirmation received")
	}
	if !confirm.Ack {
		return fmt.Errorf("broker nack message with delivery tag %d", confirm.DeliveryTag)
	}

	log.Printf(
		"published order report: branch_id=%d, total=%.2f",
		report.BranchID,
		report.Total,
	)

	return nil
}
