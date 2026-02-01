package models

type Product struct {
	ID       int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string  `gorm:"type:varchar(100);not null" json:"name"`
	Price    float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQty int64   `gorm:"not null" json:"stock_qty"`
}

func (Product) TableName() string {
	return "products"
}
