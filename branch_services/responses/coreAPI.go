package responses

type CoreAPI struct {
	Status bool           `json:"status"`
	Code   int            `json:"code"`
	Data   []StoreProduct `json:"data"`
}
