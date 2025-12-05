package invoice

import (
	"context"
	"fmt"

	"Veritasbackend/internal/domain/repositories"
)

type CreateInvoiceUseCase struct {
	invoiceRepo repositories.InvoiceRepository
	productRepo  repositories.ProductRepository
}

func NewCreateInvoiceUseCase(invoiceRepo repositories.InvoiceRepository, productRepo repositories.ProductRepository) *CreateInvoiceUseCase {
	return &CreateInvoiceUseCase{
		invoiceRepo: invoiceRepo,
		productRepo: productRepo,
	}
}

type InvoiceItemRequest struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type CreateInvoiceRequest struct {
	Items []InvoiceItemRequest `json:"items"`
}

type InvoiceItemDTO struct {
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	Subtotal  float64 `json:"subtotal"`
	ProductName string `json:"productName"`
}

type InvoiceDTO struct {
	ID        int             `json:"id"`
	Total     float64         `json:"total"`
	Status    string          `json:"status"`
	UserID    int             `json:"userId"`
	Items     []InvoiceItemDTO `json:"items"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

func (uc *CreateInvoiceUseCase) Execute(ctx context.Context, tenantID, userID int, req CreateInvoiceRequest) (*InvoiceDTO, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("la factura debe tener al menos un item")
	}

	// Validar productos y calcular totales
	var total float64
	repoItems := make([]repositories.InvoiceItem, 0, len(req.Items))
	itemDTOs := make([]InvoiceItemDTO, 0, len(req.Items))

	for _, item := range req.Items {
		// Obtener producto
		product, err := uc.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("producto con ID %d no encontrado", item.ProductID)
		}

		// Validar que el producto pertenece al tenant
		if product.TenantID != tenantID {
			return nil, fmt.Errorf("producto con ID %d no pertenece a tu tenant", item.ProductID)
		}

		// Validar stock
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("stock insuficiente para producto %s: disponible %d, solicitado %d", product.Name, product.Stock, item.Quantity)
		}

		// Validar cantidad
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("la cantidad debe ser mayor a 0")
		}

		// Calcular subtotal
		unitPrice := product.Price
		subtotal := unitPrice * float64(item.Quantity)
		total += subtotal

		// Actualizar stock del producto
		err = uc.productRepo.UpdateStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error al actualizar stock: %v", err)
		}

		// Preparar item para el repositorio
		repoItems = append(repoItems, repositories.InvoiceItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: unitPrice,
			Subtotal:  subtotal,
		})

		// Preparar DTO
		itemDTOs = append(itemDTOs, InvoiceItemDTO{
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  unitPrice,
			Subtotal:   subtotal,
			ProductName: product.Name,
		})
	}

	// Crear factura
	inv, err := uc.invoiceRepo.Create(ctx, tenantID, userID, total, repoItems)
	if err != nil {
		return nil, fmt.Errorf("error al crear factura: %v", err)
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

