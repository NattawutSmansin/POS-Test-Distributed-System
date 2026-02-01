package consumer

import (
	"cor-service/models"
	"cor-service/requests"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateFromBranch(request requests.OrderReport) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		order := models.Order{
			BranchID:  int64(request.BranchID),
			ProductID: int64(request.ProductID),
			Qty:       int64(request.Qty),
			Price:     request.Price,
			Total:     request.Total,
		}

		if err := r.CutStockProduct(tx, request); err != nil {
			return err
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *OrderRepository) CutStockProduct(tx *gorm.DB, request requests.OrderReport) error {
	baseSql := `
		UPDATE store_products
		SET qty = qty - ?
		WHERE product_id = ?
		AND branch_id = ?
	`

	result := tx.Exec(
		baseSql,
		request.Qty,
		request.ProductID,
		request.BranchID,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
