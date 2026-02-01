package domains

import "branch-service/responses"

type UseCase interface {
	BranchProductList(branchID string) (response []responses.StoreProduct, err error)
}

type Repository interface {}
