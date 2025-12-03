package dashboard

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
)

type GetMetricsUseCase struct {
	productRepo repositories.ProductRepository
	invoiceRepo repositories.InvoiceRepository
}

func NewGetMetricsUseCase(productRepo repositories.ProductRepository, invoiceRepo repositories.InvoiceRepository) *GetMetricsUseCase {
	return &GetMetricsUseCase{
		productRepo: productRepo,
		invoiceRepo: invoiceRepo,
	}
}

type MetricsResponse struct {
	TotalProducts int     `json:"totalProducts"`
	TotalInvoices int     `json:"totalInvoices"`
	Revenue       float64 `json:"revenue"`
	LowStockItems int     `json:"lowStockItems"`
}

func (uc *GetMetricsUseCase) Execute(ctx context.Context, tenantID int) (*MetricsResponse, error) {
	totalProducts, err := uc.productRepo.CountByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	totalInvoices, err := uc.invoiceRepo.CountByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Por ahora, revenue es 0 hasta que tengamos facturas reales
	revenue := 0.0

	// Low stock items (menos de 10 unidades)
	// Por simplicidad, retornamos 0. Se puede implementar después
	lowStockItems := 0

	return &MetricsResponse{
		TotalProducts: totalProducts,
		TotalInvoices: totalInvoices,
		Revenue:       revenue,
		LowStockItems: lowStockItems,
	}, nil
}

