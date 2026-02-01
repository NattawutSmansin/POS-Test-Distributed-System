package main

import (
	"fmt"
	"log"
	"time"

	DB "cor-service/databases"
	"cor-service/routes"
	rabbitmq "cor-service/services"
	"cor-service/services/consumer"
	"cor-service/services/publisher"

	"github.com/joho/godotenv"
)

// @title POS Test
// @version 1.0
// @description Example Core Service
// @host localhost:3344
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error load .env file")
	}

	DB.ConnectDatabase()

	initTimeZone()

	r := routes.SetupRouter()

	connRabbitmq()

	r.Run(":3344")
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

	fmt.Print("Rabbitmq Connect")

	consumer := consumer.NewOrderReportConsumer(conn, DB.DB)
	publisher := publisher.NewPriceUpdatePublisher(conn)
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
