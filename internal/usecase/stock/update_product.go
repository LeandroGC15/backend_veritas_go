package stock

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
	pkg_errors "Veritasbackend/pkg/errors"
)

type UpdateProductUseCase struct {
	productRepo repositories.ProductRepository
}

func NewUpdateProductUseCase(productRepo repositories.ProductRepository) *UpdateProductUseCase {
	return &UpdateProductUseCase{
		productRepo: productRepo,
	}
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	SKU         string  `json:"sku"`
}

func (uc *UpdateProductUseCase) Execute(ctx context.Context, id int, req UpdateProductRequest) (*ProductDTO, error) {
	product, err := uc.productRepo.Update(ctx, id, req.Name, req.Description, req.SKU, req.Price, req.Stock)
	if err != nil {
		return nil, pkg_errors.ErrNotFound
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

