package usecases

import (
	"cor-service/responses"
	"cor-service/src/order/domains"
)

type OrderUseCase struct {
	orderRepo domains.Repository
}

func NewOrderUseCase(repo domains.Repository) domains.UseCase {
	return &OrderUseCase{
		orderRepo: repo,
	}
}

func (u *OrderUseCase) SalesReport() (response responses.BranchSalesReport, err error) {
	saleReports, err := u.orderRepo.SalesReport()
	if err != nil {
		return response, err
	}

	if len(saleReports) == 0 {
		saleReports = []responses.SalesReportRes{}
	}

	response.BranchReports = saleReports
	return response, nil
}
