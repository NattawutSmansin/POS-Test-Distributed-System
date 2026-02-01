package domains

import (
	"cor-service/responses"
)

type UseCase interface {
	SalesReport() (response responses.BranchSalesReport, err error)
}

type Repository interface {
	SalesReport() (response []responses.SalesReportRes, err error)
}
