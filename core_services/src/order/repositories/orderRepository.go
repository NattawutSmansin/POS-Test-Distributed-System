package repositories

import (
	"cor-service/responses"
	"cor-service/src/order/handlers"
	"cor-service/src/order/usecases"

	"gorm.io/gorm"
)

type OrderRepository struct {
	conn *gorm.DB
}

func NewRepositoryHandler(conn *gorm.DB) *handlers.OrderHandler {
	useCase := usecases.NewOrderUseCase(&OrderRepository{conn})
	handler := handlers.NewOrderHandler(useCase)
	return handler
}

func (r *OrderRepository) SalesReport() (response []responses.SalesReportRes, err error) {
	sqlBase := `
		SELECT 
			branches.name AS branch_name,
			COUNT(orders.id) AS order_count,
			COALESCE(SUM(orders.price), 0) AS order_price,
			COALESCE(
				SUM(orders.price) / NULLIF(COUNT(orders.id), 0),
				0
			) AS order_avg_price,
			COALESCE(SUM(orders.qty), 0) AS order_qty 
		FROM branches
		LEFT JOIN orders 
		ON branches.id  = orders.branch_id
		GROUP BY 
			branches.id,
			branches.name;
	`

	var results []responses.SalesReportRes
	if err := r.conn.Raw(sqlBase).Scan(&results).Error; err != nil {
		return response, err
	}

	return results, nil
}
