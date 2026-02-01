package requests

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	StockQty int64   `json:"stock_qty"`
}

type ProductPrice struct {
	ID    int64   `json:"id"`
	Price float64 `json:"price"`
}
