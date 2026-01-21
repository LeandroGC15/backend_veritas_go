package supplier

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
)

type ListSuppliersUseCase struct {
	supplierRepo repositories.SupplierRepository
}

func NewListSuppliersUseCase(supplierRepo repositories.SupplierRepository) *ListSuppliersUseCase {
	return &ListSuppliersUseCase{
		supplierRepo: supplierRepo,
	}
}

type ListSuppliersRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type SuppliersListResponse struct {
	Suppliers []SupplierDTO `json:"suppliers"`
	Total     int           `json:"total"`
	Page      int           `json:"page"`
	Limit     int           `json:"limit"`
}

func (uc *ListSuppliersUseCase) Execute(ctx context.Context, tenantID int, req ListSuppliersRequest) (*SuppliersListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	offset := (req.Page - 1) * req.Limit

	suppliers, total, err := uc.supplierRepo.FindAll(ctx, tenantID, req.Limit, offset)
	if err != nil {
		return nil, err
	}

	supplierDTOs := make([]SupplierDTO, len(suppliers))
	for i, supplier := range suppliers {
		supplierDTOs[i] = *convertSupplierToDTO(supplier)
	}

	return &SuppliersListResponse{
		Suppliers: supplierDTOs,
		Total:     total,
		Page:      req.Page,
		Limit:     req.Limit,
	}, nil
}