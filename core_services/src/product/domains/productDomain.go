package domains

import (
	"cor-service/requests"
	"cor-service/responses"
)

type UseCase interface {
	CreateProduct(request requests.Product) (err error)
	UpdateProductPrice(request requests.ProductPrice) (err error)
	BranchProductList(branchID int64) (response []responses.StoreProduct, err error)
}

type Repository interface {
	Create(model interface{}) (err error)
	UpdateProductPrice(request requests.ProductPrice) (err error)
	BranchProductList(branchID int64) (response []responses.StoreProduct, err error)
}
