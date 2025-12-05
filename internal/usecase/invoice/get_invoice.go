package invoice

import (
	"context"
	"fmt"

	"Veritasbackend/internal/domain/repositories"
)

type GetInvoiceUseCase struct {
	invoiceRepo repositories.InvoiceRepository
	productRepo repositories.ProductRepository
}

func NewGetInvoiceUseCase(invoiceRepo repositories.InvoiceRepository, productRepo repositories.ProductRepository) *GetInvoiceUseCase {
	return &GetInvoiceUseCase{
		invoiceRepo: invoiceRepo,
		productRepo: productRepo,
	}
}

func (uc *GetInvoiceUseCase) Execute(ctx context.Context, invoiceID int) (*InvoiceDTO, error) {
	inv, items, err := uc.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("factura no encontrada: %v", err)
	}

	itemDTOs := make([]InvoiceItemDTO, len(items))
	for i, item := range items {
		// Obtener información del producto
		product, err := uc.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			// Si no se encuentra el producto, usar valores por defecto
			product = nil
		}

		productName := ""
		if product != nil {
			productName = product.Name
		}

		itemDTOs[i] = InvoiceItemDTO{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			Subtotal:   item.Subtotal,
			ProductName: productName,
		}
	}

	return &InvoiceDTO{
		ID:        inv.ID,
		Total:     inv.Total,
		Status:    inv.Status,
		UserID:    inv.UserID,
		Items:     itemDTOs,
		CreatedAt: inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: inv.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

