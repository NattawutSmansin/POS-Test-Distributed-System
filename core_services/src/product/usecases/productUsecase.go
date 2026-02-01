package usecases

import (
	"cor-service/models"
	"cor-service/requests"
	"cor-service/responses"
	"cor-service/src/product/domains"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ProductUseCase struct {
	productRepo domains.Repository
}

func NewProductUseCase(repo domains.Repository) domains.UseCase {
	return &ProductUseCase{
		productRepo: repo,
	}
}

func (u *ProductUseCase) CreateProduct(request requests.Product) (err error) {
	mapProductModel := models.Product{
		Name:     request.Name,
		Price:    request.Price,
		StockQty: request.StockQty,
	}

	err = u.productRepo.Create(&mapProductModel)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductUseCase) BranchProductList(branchID int64) (response []responses.StoreProduct, err error) {
	response, err = u.productRepo.BranchProductList(branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, fmt.Errorf("not found for ID %d", branchID)
		}
		return response, err
	}
	return response, nil
}

func (u *ProductUseCase) UpdateProductPrice(request requests.ProductPrice) (err error) {
	err = u.productRepo.UpdateProductPrice(request)
	if err != nil {
		return err
	}
	return nil
}
