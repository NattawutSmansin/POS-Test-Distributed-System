package handlers

import (
	"branch-service/helpers"
	"branch-service/requests"
	"branch-service/responses"
	"branch-service/src/order/domains"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase domains.UseCase
}

func NewOrderHandler(usecase domains.UseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: usecase,
	}
}

// @summary
// @description 2.2 POST /orders – ทำรายการขายสินค้า
// @tags Branch Services
// @id CreateOrder
// @param body body requests.OrderReport true "body"
// @response 201 {object} responses.Success{} "OK"
// @response 409 "Conflict Error"
// @response 422 "Entity Error"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/branch/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request requests.OrderReport
	if err := c.ShouldBindJSON(&request); err != nil {
		newError := responses.NewAppErr(http.StatusBadRequest, err.Error())
		errResponse := responses.FailRespone(newError)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	validate, err := ValidateOrder(request)
	if err != nil {
		newError := responses.NewAppErr(http.StatusBadRequest, err.Error())
		errResponse := responses.FailRespone(newError)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if len(validate) > 0 {
		errResponse := responses.ValidateResponse(validate)
		c.JSON(http.StatusUnprocessableEntity, errResponse)
		return
	}

	err = h.orderUseCase.CreateOrder(request)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, responses.SuccessResponse(nil, http.StatusCreated))
}

func ValidateOrder(request requests.OrderReport) (map[string]string, error) {
	fiedValidate := make(map[string]interface{})

	fiedValidate["name"] = request.BranchID
	fiedValidate["name"] = request.ProductID
	fiedValidate["price"] = request.Price
	fiedValidate["qty"] = request.Qty
	fiedValidate["total"] = request.Total

	makeValidate, err := helpers.MakeValidate(fiedValidate)
	if err != nil {
		return nil, err
	}

	return makeValidate, nil
}
