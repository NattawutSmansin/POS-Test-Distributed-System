package domains

import "branch-service/requests"

type UseCase interface {
	CreateOrder(request requests.OrderReport) (err error)
}

type Repository interface {
	CreateOrder(request requests.OrderReport) (err error)
}
