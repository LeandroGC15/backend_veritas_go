package stock

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
	pkg_errors "Veritasbackend/pkg/errors"
)

type DeleteProductUseCase struct {
	productRepo repositories.ProductRepository
}

func NewDeleteProductUseCase(productRepo repositories.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{
		productRepo: productRepo,
	}
}

func (uc *DeleteProductUseCase) Execute(ctx context.Context, id int) error {
	err := uc.productRepo.Delete(ctx, id)
	if err != nil {
		return pkg_errors.ErrNotFound
	}
	return nil
}

