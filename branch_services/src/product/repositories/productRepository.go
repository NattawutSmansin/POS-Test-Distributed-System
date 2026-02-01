package repositories

import (
	"branch-service/src/product/handlers"
	"branch-service/src/product/usecases"
)

type ProductRepository struct{}

func NewRepositoryHandler() *handlers.ProductHandler {
	useCase := usecases.NewProductUseCase(&ProductRepository{})
	handler := handlers.NewProductHandler(useCase)
	return handler
}
