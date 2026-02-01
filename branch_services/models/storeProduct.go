package models

type StoreProduct struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID  int64   `gorm:"not null" json:"branch_id"`
	ProductID int64   `gorm:"not null" json:"product_id"`
	Qty       int64   `gorm:"not null" json:"qty"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive  bool    `gorm:"not null;default:false" json:"is_active"`

	// Relations
	Branchs  []Branch  `gorm:"foreignKey:BranchID;references:ID" json:"branchs"`
	Products []Product `gorm:"foreignKey:ProductID;references:ID" json:"products"`
}

func (StoreProduct) TableName() string {
	return "store_product"
}
