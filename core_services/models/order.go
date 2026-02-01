package models

type Order struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID  int64   `gorm:"not null" json:"branch_id"`
	ProductID int64   `gorm:"not null" json:"product_id"`
	Qty       int64   `gorm:"not null" json:"qty"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Total     float64 `gorm:"type:decimal(10,2);not null" json:"total"`
}

func (Order) TableName() string {
	return "orders"
}