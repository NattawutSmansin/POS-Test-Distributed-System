package usecases

import (
	"branch-service/responses"
	"branch-service/src/product/domains"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ProductUseCase struct {
	productRepo domains.Repository
}

func NewProductUseCase(repo domains.Repository) domains.UseCase {
	return &ProductUseCase{
		productRepo: repo,
	}
}

func (u *ProductUseCase) BranchProductList(branchID string) (response []responses.StoreProduct, err error) {
	apiURL := os.Getenv("COR_URL")
	pathURL := fmt.Sprintf("%s/%s/products", apiURL, branchID)

	req, err := http.NewRequest(http.MethodGet, pathURL, nil)
	if err != nil {
		return response, responses.NewAppErr(400, err.Error())
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return response, responses.NewAppErr(500, err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return response, responses.NewAppErr(
			res.StatusCode,
			"Core service error: "+string(body),
		)
	}

	var apiResp responses.CoreAPI

	if err := json.NewDecoder(res.Body).Decode(&apiResp); err != nil {
		return response, responses.NewAppErr(500, err.Error())
	}

	return apiResp.Data, nil
}
