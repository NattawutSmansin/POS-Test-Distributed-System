package responses

type Validate struct {
	Status bool `json:"status" exanple:"false" extensions:"x-order=0"`
	Code   int  `json:"code" example:"422" extensions:"x-order=1"`
	Err    struct {
		Message string            `json:"message" example:"validate error" extensions:"x-order=0"`
		Field   map[string]string `json:"field" extensions:"x-order=1"`
	} `json:"error" extensions:"x-order=2"`
}

type InternalServerErrorResponse struct {
	Status bool `json:"status" example:"false" extensions:"x-order=0"`
	Code   int  `json:"code" example:"500" extensions:"x-order=1"`
	Error  struct {
		Message string   `json:"message" example:"internal server error" extensions:"x-order=0"`
		Field   struct{} `json:"field" extensions:"x-order=1"`
	} `json:"error" extensions:"x-order=2"`
}

type UnauthorizedErrorResponse struct {
	Status bool `json:"status" example:"false" extensions:"x-order=0"`
	Code   int  `json:"code" example:"401" extensions:"x-order=1"`
	Error  struct {
		Message string   `json:"message" example:"unauthorized to access the resource" extensions:"x-order=0"`
		Field   struct{} `json:"field" extensions:"x-order=1"`
	} `json:"error" extensions:"x-order=2"`
}

type BadRequestErrorResponse struct {
	Status bool `json:"status" example:"false" extensions:"x-order=0"`
	Code   int  `json:"code" example:"400" extensions:"x-order=1"`
	Error  struct {
		Message string   `json:"message" example:"bad request error message" extensions:"x-order=0"`
		Field   struct{} `json:"field" extensions:"x-order=1"`
	} `json:"error" extensions:"x-order=2"`
}

type NoDataResponse struct {
	Status bool   `json:"status" example:"true" extensions:"x-order=0"`
	Code   int    `json:"code" extensions:"x-order=1"`
	Data   NoData `json:"data" extensions:"x-order=2"`
}

type CreateResponse struct {
	Status bool   `json:"status" example:"true" extensions:"x-order=0"`
	Code   int    `json:"code" extensions:"x-order=1" example:"201"`
	Data   NoData `json:"data" extensions:"x-order=2"`
}

type UpdateResponse struct {
	Status bool   `json:"status" example:"true" extensions:"x-order=0"`
	Code   int    `json:"code" extensions:"x-order=1" example:"202"`
	Data   NoData `json:"data" extensions:"x-order=2"`
}

type NotFoundErrorResponse struct {
	Status bool `json:"status" example:"false" extensions:"x-order=0"`
	Code   int  `json:"code" example:"400" extensions:"x-order=1"`
	Error  struct {
		Message string   `json:"message" example:"data was not found" extensions:"x-order=0"`
		Field   struct{} `json:"field" extensions:"x-order=1"`
	} `json:"error" extensions:"x-order=2"`
}

type Empty struct {
}
