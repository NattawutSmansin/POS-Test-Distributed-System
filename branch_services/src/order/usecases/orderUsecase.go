package usecases

import (
	"branch-service/requests"
	"branch-service/src/order/domains"
)

type OrderUseCase struct {
	orderRepo domains.Repository
}

func NewOrderUseCase(repo domains.Repository) domains.UseCase {
	return &OrderUseCase{
		orderRepo: repo,
	}
}

func (u *OrderUseCase) CreateOrder(request requests.OrderReport) (err error) {
	err = u.orderRepo.CreateOrder(request)
	if err != nil {
		return err
	}
	return nil
}
