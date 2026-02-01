package repositories

import (
	"branch-service/requests"
	rabbitmq "branch-service/services"
	orderPublisher "branch-service/services/publisher"
	"branch-service/src/order/handlers"
	"branch-service/src/order/usecases"
	"log"
)

type OrderRepository struct {
	connRabbitmq *rabbitmq.Connection
}

func NewRepositoryHandler(connRabbitmq *rabbitmq.Connection) *handlers.OrderHandler {
	useCase := usecases.NewOrderUseCase(&OrderRepository{connRabbitmq})
	handler := handlers.NewOrderHandler(useCase)
	return handler
}

func (r OrderRepository) CreateOrder(request requests.OrderReport) (err error) {
	conn := r.connRabbitmq
	publisher := orderPublisher.NewOrderReportPublisher(conn)

	if err := publisher.Setup(); err != nil {
		log.Fatal(err)
	}

	err = publisher.Publisher(request)
	if err != nil {
		return err
	}
	return nil
}
