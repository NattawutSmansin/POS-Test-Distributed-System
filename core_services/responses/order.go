package responses

type BranchSalesReport struct {
	BranchReports []SalesReportRes `json:"branch_reports"`
}

type SalesReportRes struct {
	BranchName    string  `json:"branch_name"`
	OrderCount    int32   `json:"order_count"`
	OrderPrice    float64 `json:"order_price"`
	OrderAvgPrice float64 `json:"order_avg_price"`
	OrderQty      int32   `json:"order_qty"`
}
