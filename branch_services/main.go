package main

import (
	"log"
	"time"

	"branch-service/routes"

	"github.com/joho/godotenv"

	rabbitmq "branch-service/services"
	"branch-service/services/consumer"
	"branch-service/services/publisher"
)

// @title POS Test
// @version 1.0
// @description Example Branch Service
// @host localhost:5566
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error load .env file")
	}

	initTimeZone()

	r := routes.SetupRouter()

	connRabbitmq()

	r.Run(":5566")
}

func initTimeZone() {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = location
}

func connRabbitmq() {
	conn, err := rabbitmq.NewConnectionRabbitMQ()
	if err != nil {
		panic(err)
	}

	consumer := consumer.NewPriceUpdateConsumer(conn)
	publisher := publisher.NewOrderReportPublisher(conn)
	if err := consumer.Setup(); err != nil {
		log.Fatal(err)
	}

	if err := publisher.Setup(); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := consumer.Start(); err != nil {
			log.Fatal(err)
		}
	}()
}
