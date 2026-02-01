package responses

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AppErr struct {
	StatusCode int
	Message    string
	Field      map[string]string
}

type NoData struct{}

type ErrorFormat struct {
	Message string      `json:"message" extensions:"x-order=0"`
	Field   interface{} `json:"field" extensions:"x-order=1"`
}

func (r *AppErr) Error() string {
	return r.Message
}

// 400, 500
type Fail struct {
	Status bool        `json:"status" extensions:"x-order=0"`
	Code   int         `json:"code" extensions:"x-order=1"`
	Err    ErrorFormat `json:"error" extensions:"x-order=2"`
}

// 200, 201, 202
type Success struct {
	Status bool        `json:"status" example:"true" extensions:"x-order=0"`
	Code   int         `json:"code" extensions:"x-order=1"`
	Data   interface{} `json:"data" extensions:"x-order=2"`
}

func SuccessResponse(data interface{}, code int) Success {
	return Success{
		Status: true,
		Code:   code,
		Data:   data,
	}
}

func ValidateResponse(msg map[string]string) Fail {
	var request Fail
	request.Code = 422
	request.Err.Message = "Validate error"
	request.Err.Field = msg
	request.Status = false
	return request
}

func FailRespone(err error) Fail {
	if _, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string]string)

		return Fail{
			Status: false,
			Code:   http.StatusUnprocessableEntity,
			Err: ErrorFormat{
				Message: "Validate error",
				Field:   fieldErrors,
			},
		}
	}

	// Check if it's an AppErr type
	appErr, ok := err.(*AppErr)
	if ok {
		return Fail{
			Status: false,
			Code:   appErr.StatusCode,
			Err: ErrorFormat{
				Message: appErr.Message,
				Field:   appErr.Field,
			},
		}
	}

	// Error อื่นๆ
	return Fail{
		Status: false,
		Code:   http.StatusBadRequest,
		Err: ErrorFormat{
			Message: err.Error(),
			Field:   NoData{},
		},
	}
}

func NewAppErr(statusCode int, msg string) *AppErr {
	return &AppErr{
		StatusCode: statusCode,
		Message:    msg,
	}
}
