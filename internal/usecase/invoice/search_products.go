package invoice

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
)

type SearchProductsUseCase struct {
	invoiceRepo repositories.InvoiceRepository
}

func NewSearchProductsUseCase(invoiceRepo repositories.InvoiceRepository) *SearchProductsUseCase {
	return &SearchProductsUseCase{
		invoiceRepo: invoiceRepo,
	}
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

type SearchProductsResponse struct {
	Products []ProductDTO `json:"products"`
}

func (uc *SearchProductsUseCase) Execute(ctx context.Context, tenantID int, query string) (*SearchProductsResponse, error) {
	if query == "" {
		return &SearchProductsResponse{Products: []ProductDTO{}}, nil
	}

	products, err := uc.invoiceRepo.SearchProducts(ctx, tenantID, query)
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

	return &SearchProductsResponse{
		Products: productDTOs,
	}, nil
}

