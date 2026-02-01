package repositories

import (
	"cor-service/requests"
	"cor-service/responses"
	"cor-service/src/product/handlers"
	"cor-service/src/product/usecases"
	"log"

	rabbitmq "cor-service/services"
	productPublisher "cor-service/services/publisher"

	"gorm.io/gorm"
)

type ProductRepository struct {
	conn         *gorm.DB
	connRabbitmq *rabbitmq.Connection
}

func NewRepositoryHandler(conn *gorm.DB, connRabbitmq *rabbitmq.Connection) *handlers.ProductHandler {
	useCase := usecases.NewProductUseCase(&ProductRepository{conn, connRabbitmq})
	handler := handlers.NewProductHandler(useCase)
	return handler
}

func (r *ProductRepository) Create(model interface{}) (err error) {
	query := r.conn.Begin()
	query.Debug()
	if err = query.Create(model).Error; err != nil {
		query.Rollback()
		return err
	}
	query.Commit()
	return nil
}

func (r *ProductRepository) BranchProductList(branchID int64) (response []responses.StoreProduct, err error) {
	sqlBase := `
		SELECT 
			products.id AS product_id, 
			products.name, 
			store_products.price, 
			store_products.qty
		FROM branches 
		JOIN store_products 
			ON store_products.branch_id = branches.id
			AND store_products.is_active = true
		JOIN products 
			ON store_products.product_id = products.id
		WHERE 
			branches.is_active = true
		AND branches.id = ?
	`

	var products []responses.StoreProduct

	if err := r.conn.Raw(sqlBase, branchID).Scan(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) UpdateProductPrice(request requests.ProductPrice) (err error) {
	if err := r.conn.Transaction(func(tx *gorm.DB) error {
		baseSql := `
			UPDATE products
			SET price = ?
			WHERE id = ?
		`

		result := tx.Exec(baseSql, request.Price, request.ID)
		if result.Error != nil {
			return result.Error
		}

		if err := r.UpdateBranchProductPrice(tx, request); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	publisher := productPublisher.NewPriceUpdatePublisher(r.connRabbitmq)
	go func() {
		if err := publisher.Setup(); err != nil {
			log.Fatal(err)
		}

		if err := publisher.Publish(request); err != nil {
			log.Println("publish price update failed:", err)
		}
	}()

	return nil
}

func (r *ProductRepository) UpdateBranchProductPrice(tx *gorm.DB, request requests.ProductPrice) error {
	baseSql := `
		UPDATE store_products
		SET price = ?
		WHERE product_id = ?
	`

	result := tx.Exec(
		baseSql,
		request.Price,
		request.ID,
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
