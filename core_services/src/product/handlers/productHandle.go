package handlers

import (
	"cor-service/helpers"
	"cor-service/requests"
	"cor-service/responses"
	"cor-service/src/product/domains"
	"fmt"
	"net/http"
	"strconv"

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
// @description 1.1 POST /products – สร้างสินค้าใหม่
// @tags Core Services
// @id CreateProduct
// @param body body requests.Product true "body"
// @response 201 {object} responses.Success{} "OK"
// @response 409 "Conflict Error"
// @response 422 "Entity Error"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/core/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request requests.Product
	if err := c.ShouldBindJSON(&request); err != nil {
		newError := responses.NewAppErr(http.StatusBadRequest, err.Error())
		errResponse := responses.FailRespone(newError)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	validate, err := ValidateProduct(request)
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

	err = h.productUseCase.CreateProduct(request)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, responses.SuccessResponse(nil, http.StatusCreated))
}

func ValidateProduct(request requests.Product) (map[string]string, error) {
	fiedValidate := make(map[string]interface{})

	fiedValidate["name"] = request.Name
	fiedValidate["price"] = request.Price
	fiedValidate["stock_qty"] = request.StockQty

	makeValidate, err := helpers.MakeValidate(fiedValidate)
	if err != nil {
		return nil, err
	}

	return makeValidate, nil
}

// @summary
// @description 2.1 GET /products – แสดงข้อมูลสินค้าทั้งหมดของสาขา
// @tags Core Services
// @id BranchProduct
// @Param branch_id path string true "ID สาขา"
// @response 200 {object} responses.Success{} "OK"
// @response 409 "Conflict Error"
// @response 422 "Entity Error"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/core/{branch_id}/products [get]
func (h *ProductHandler) BranchProduct(c *gin.Context) {
	var branchID string = c.Param("branch_id")
	branchIDConvert, _ := strconv.ParseInt(fmt.Sprintf("%v", branchID), 10, 64)

	branchProductList, err := h.productUseCase.BranchProductList(branchIDConvert)
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

// @summary
// @description 3.1 การส่งข้อมูลจาก Core ไปยัง Branch การอัปเดตราคาสินค้าให้ทุกสาขาพร้อมกัน
// @tags Core Services
// @id UpdateProductPrice
// @param body body requests.ProductPrice true "body"
// @response 202 {object} responses.Success{} "OK"
// @response 409 "Conflict Error"
// @response 422 "Entity Error"
// @response 500 {object} responses.InternalServerErrorResponse "Internal Server Error"
// @router /api/core/products [put]
func (h *ProductHandler) UpdateProductPrice(c *gin.Context) {
	var request requests.ProductPrice
	if err := c.ShouldBindJSON(&request); err != nil {
		newError := responses.NewAppErr(http.StatusBadRequest, err.Error())
		errResponse := responses.FailRespone(newError)
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	validate, err := ValidateUpdateProductPrice(request)
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

	err = h.productUseCase.UpdateProductPrice(request)
	if err != nil {
		errResponse := responses.FailRespone(err)
		c.AbortWithStatusJSON(errResponse.Code, errResponse)
		return
	}

	c.AbortWithStatusJSON(http.StatusAccepted, responses.SuccessResponse(nil, http.StatusAccepted))
}

func ValidateUpdateProductPrice(request requests.ProductPrice) (map[string]string, error) {
	fiedValidate := make(map[string]interface{})
	fiedValidate["price"] = request.Price

	makeValidate, err := helpers.MakeValidate(fiedValidate)
	if err != nil {
		return nil, err
	}

	return makeValidate, nil
}
