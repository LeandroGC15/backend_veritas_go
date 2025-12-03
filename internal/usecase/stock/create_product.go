package stock

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
)

type CreateProductUseCase struct {
	productRepo repositories.ProductRepository
}

func NewCreateProductUseCase(productRepo repositories.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepo: productRepo,
	}
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	SKU         string  `json:"sku"`
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, tenantID int, req CreateProductRequest) (*ProductDTO, error) {
	product, err := uc.productRepo.Create(ctx, tenantID, req.Name, req.Description, req.SKU, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}

	return &ProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.Sku,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

