package stock

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
)

type ListProductsUseCase struct {
	productRepo repositories.ProductRepository
}

func NewListProductsUseCase(productRepo repositories.ProductRepository) *ListProductsUseCase {
	return &ListProductsUseCase{
		productRepo: productRepo,
	}
}

type ListProductsRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ProductDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	SKU         string  `json:"sku"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type ListProductsResponse struct {
	Products []ProductDTO `json:"products"`
	Total    int          `json:"total"`
	Page     int          `json:"page"`
	Limit    int          `json:"limit"`
}

func (uc *ListProductsUseCase) Execute(ctx context.Context, tenantID int, req ListProductsRequest) (*ListProductsResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	offset := (req.Page - 1) * req.Limit

	products, total, err := uc.productRepo.FindAll(ctx, tenantID, req.Limit, offset)
	if err != nil {
		return nil, err
	}

	productDTOs := make([]ProductDTO, len(products))
	for i, p := range products {
		productDTOs[i] = ProductDTO{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			SKU:         p.Sku,
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &ListProductsResponse{
		Products: productDTOs,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

