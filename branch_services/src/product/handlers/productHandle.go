package handlers

import (
	"branch-service/responses"
	"branch-service/src/product/domains"

	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUseCase domains.UseCase
}

func NewProductHandler(usecase domains.UseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: usecase,
	}
}

// @summary
// @description 2.1 GET /products – แสดงข้อมูลสินค้าทั้งหมดของสาขา
// @tags Branch Services
// @id BranchProduct
// @Param branch_id path string true "ID สาขา"
// @response 200 {object} responses.Success{} "OK"
// @response 409 "Conflict Error"
// @response 422 "Entity Error"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/branch/{branch_id}/products [get]
func (h *ProductHandler) BranchProduct(c *gin.Context) {
	var branchID string = c.Param("branch_id")

	branchProductList, err := h.productUseCase.BranchProductList(branchID)
	if err != nil {
		newError := responses.NewAppErr(http.StatusInternalServerError, err.Error())
		errResponse := responses.FailRespone(newError)
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	if len(branchProductList) == 0 {
		branchProductList = []responses.StoreProduct{}
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(branchProductList, http.StatusOK))
}
