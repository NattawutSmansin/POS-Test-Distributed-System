package requests

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	StockQty int64   `json:"stock_qty"`
}
