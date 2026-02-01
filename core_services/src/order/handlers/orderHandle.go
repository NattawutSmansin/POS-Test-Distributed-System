package handlers

import (
	"cor-service/responses"
	"cor-service/src/order/domains"
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
// @description 1.2 GET /orders – รายงานยอดขาย
// @tags Core Services
// @id Order
// @response 200 {object} responses.Success{} "OK"
// @response 409 "Conflict Error"
// @response 422 "Entity Error"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/core/orders [get]
func (h *OrderHandler) Order(c *gin.Context) {
	response, err := h.orderUseCase.SalesReport()
	if err != nil {
		newError := responses.NewAppErr(http.StatusInternalServerError, err.Error())
		errResponse := responses.FailRespone(newError)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, responses.SuccessResponse(response, http.StatusOK))
}
