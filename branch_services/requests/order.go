package requests

type OrderReport struct {
	BranchID  int     `json:"branch_id"`
	ProductID int     `json:"product_id"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
	Total     float64 `json:"total"`
}