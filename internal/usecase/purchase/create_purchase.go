package purchase

import (
	"context"
	"fmt"
	"log"
	"time"

	"Veritasbackend/ent"
	"Veritasbackend/internal/domain/repositories"
)

type CreatePurchaseUseCase struct {
	purchaseInvoiceRepo     repositories.PurchaseInvoiceRepository
	purchaseInvoiceItemRepo repositories.PurchaseInvoiceItemRepository
	productRepo             repositories.ProductRepository
}

func NewCreatePurchaseUseCase(
	purchaseInvoiceRepo repositories.PurchaseInvoiceRepository,
	purchaseInvoiceItemRepo repositories.PurchaseInvoiceItemRepository,
	productRepo repositories.ProductRepository,
) *CreatePurchaseUseCase {
	return &CreatePurchaseUseCase{
		purchaseInvoiceRepo:     purchaseInvoiceRepo,
		purchaseInvoiceItemRepo: purchaseInvoiceItemRepo,
		productRepo:             productRepo,
	}
}

type PurchaseItemRequest struct {
	ProductID    int     `json:"productId"`
	Quantity     int     `json:"quantity"`
	UnitCost     float64 `json:"unitCost"`
	ProductName  string  `json:"productName"`
	ProductSku   string  `json:"productSku"`
	ProductPrice float64 `json:"productPrice"`
}

type CreatePurchaseRequest struct {
	SupplierID    int                   `json:"supplierId"`
	InvoiceNumber string                `json:"invoiceNumber"`
	PaymentMethod *string               `json:"paymentMethod,omitempty"`
	DueDate       *time.Time            `json:"dueDate,omitempty"`
	Items         []PurchaseItemRequest `json:"items"`
}

type PurchaseItemDTO struct {
	ID                int     `json:"id"`
	PurchaseInvoiceID int     `json:"purchaseInvoiceId"`
	ProductID         int     `json:"productId"`
	Quantity          int     `json:"quantity"`
	UnitCost          float64 `json:"unitCost"`
	Subtotal          float64 `json:"subtotal"`
}

type PurchaseInvoiceDTO struct {
	ID            int                `json:"id"`
	InvoiceNumber string             `json:"invoiceNumber"`
	Total         float64            `json:"total"`
	Status        string             `json:"status"`
	PaymentMethod *string            `json:"paymentMethod,omitempty"`
	DueDate       *string            `json:"dueDate,omitempty"`
	PaidAmount    float64            `json:"paidAmount"`
	SupplierID    int                `json:"supplierId"`
	TenantID      int                `json:"tenantId"`
	UserID        int                `json:"userId"`
	Items         []PurchaseItemDTO  `json:"items"`
	CreatedAt     string             `json:"createdAt"`
	UpdatedAt     string             `json:"updatedAt"`
}

func convertPurchaseInvoiceToDTO(invoice *ent.PurchaseInvoice) *PurchaseInvoiceDTO {
	var dueDateStr *string
	var paymentMethod *string

	// Handle DueDate - check if it's zero time (not set)
	if !invoice.DueDate.IsZero() {
		formatted := invoice.DueDate.Format("2006-01-02")
		dueDateStr = &formatted
	}

	// Handle PaymentMethod - check if it's empty
	if invoice.PaymentMethod != "" {
		paymentMethod = &invoice.PaymentMethod
	}

	return &PurchaseInvoiceDTO{
		ID:            invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
		Total:         invoice.Total,
		Status:        invoice.Status,
		PaymentMethod: paymentMethod,
		DueDate:       dueDateStr,
		PaidAmount:    invoice.PaidAmount,
		SupplierID:    invoice.SupplierID,
		TenantID:      invoice.TenantID,
		UserID:        invoice.UserID,
		CreatedAt:     invoice.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     invoice.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func convertPurchaseInvoiceItemToDTO(item *ent.PurchaseInvoiceItem) PurchaseItemDTO {
	return PurchaseItemDTO{
		ID:                item.ID,
		PurchaseInvoiceID: item.PurchaseInvoiceID,
		ProductID:         item.ProductID,
		Quantity:          item.Quantity,
		UnitCost:          item.UnitCost,
		Subtotal:          item.Subtotal,
	}
}

func (uc *CreatePurchaseUseCase) Execute(ctx context.Context, tenantID, userID int, req CreatePurchaseRequest) (*PurchaseInvoiceDTO, error) {
	// Log the request for debugging
	log.Printf("Creating purchase for tenant %d, user %d, supplier %d, invoice %s", tenantID, userID, req.SupplierID, req.InvoiceNumber)
	log.Printf("Items in request: %d", len(req.Items))

	for i, item := range req.Items {
		log.Printf("Item %d: ProductID=%d, Quantity=%d, UnitCost=%f, ProductName='%s'", i, item.ProductID, item.Quantity, item.UnitCost, item.ProductName)
	}

	// Calculate total
	var total float64
	var purchaseItems []*ent.PurchaseInvoiceItem

	for _, item := range req.Items {
		var productID = item.ProductID

		// If ProductID is negative, it means we need to create a new product
		if item.ProductID < 0 {
			log.Printf("Creating new product for item with ProductID %d", item.ProductID)

			if item.ProductName == "" {
				log.Printf("ProductName is empty")
				return nil, fmt.Errorf("nombre del producto requerido para productos nuevos")
			}

			productName := item.ProductName
			log.Printf("Creating product with name: %s", productName)

			// Create new product
			description := ""
			sku := item.ProductSku
			price := item.UnitCost // Default price is the cost

			if sku == "" {
				sku = fmt.Sprintf("SKU-%d", item.ProductID*-1) // Generate SKU from negative ID
			}
			if item.ProductPrice > 0 {
				price = item.ProductPrice
			}

			log.Printf("Product SKU: %s, Price: %f", sku, price)

			log.Printf("Calling productRepo.Create with tenantID=%d, name=%s, sku=%s, price=%f", tenantID, productName, sku, price)

			newProduct, err := uc.productRepo.Create(ctx, tenantID, productName, description, sku, price, 0)
			if err != nil {
				log.Printf("Error creating product: %v", err)
				return nil, fmt.Errorf("error al crear nuevo producto '%s': %w", productName, err)
			}
			log.Printf("New product created with ID: %d", newProduct.ID)
			productID = newProduct.ID
		}

		subtotal := float64(item.Quantity) * item.UnitCost
		total += subtotal

		purchaseItem := &ent.PurchaseInvoiceItem{
			PurchaseInvoiceID: 0, // Will be set after invoice creation
			ProductID:         productID,
			Quantity:          item.Quantity,
			UnitCost:          item.UnitCost,
			Subtotal:          subtotal,
		}
		purchaseItems = append(purchaseItems, purchaseItem)

			// Update product stock
		log.Printf("Updating stock for product %d, quantity %d", productID, item.Quantity)
		product, err := uc.productRepo.FindByID(ctx, productID)
		if err != nil {
			log.Printf("Error finding product %d: %v", productID, err)
			return nil, fmt.Errorf("error al buscar producto con ID %d: %w", productID, err)
		}
		if product == nil {
			log.Printf("Product %d not found", productID)
			return nil, fmt.Errorf("producto con ID %d no encontrado", productID)
		}

		// Update stock: stock += quantity purchased
		err = uc.productRepo.AddStock(ctx, productID, item.Quantity)
		if err != nil {
			log.Printf("Error updating stock for product %d: %v", productID, err)
			return nil, fmt.Errorf("error al actualizar stock del producto %d: %w", productID, err)
		}
		log.Printf("Stock updated successfully for product %d", productID)
	}

	// Create purchase invoice
	log.Printf("Creating purchase invoice with total %f", total)
	invoice, err := uc.purchaseInvoiceRepo.Create(ctx, tenantID, req.SupplierID, userID, req.InvoiceNumber, total, req.PaymentMethod, req.DueDate)
	if err != nil {
		log.Printf("Error creating purchase invoice: %v", err)
		return nil, fmt.Errorf("error al crear factura de compra: %w", err)
	}
	log.Printf("Purchase invoice created with ID %d", invoice.ID)

	// Create purchase invoice items
	var itemsDTO []PurchaseItemDTO
	for _, item := range purchaseItems {
		item.PurchaseInvoiceID = invoice.ID
		log.Printf("Creating purchase invoice item for product %d", item.ProductID)
		createdItem, err := uc.purchaseInvoiceItemRepo.Create(ctx, item.PurchaseInvoiceID, item.ProductID, item.Quantity, item.UnitCost, item.Subtotal)
		if err != nil {
			log.Printf("Error creating purchase invoice item for product %d: %v", item.ProductID, err)
			return nil, fmt.Errorf("error al crear ítem de factura para producto %d: %w", item.ProductID, err)
		}

		itemsDTO = append(itemsDTO, PurchaseItemDTO{
			ID:                createdItem.ID,
			PurchaseInvoiceID: createdItem.PurchaseInvoiceID,
			ProductID:         createdItem.ProductID,
			Quantity:          createdItem.Quantity,
			UnitCost:          createdItem.UnitCost,
			Subtotal:          createdItem.Subtotal,
		})
	}
	log.Printf("Purchase completed successfully with %d items", len(itemsDTO))

	dto := convertPurchaseInvoiceToDTO(invoice)
	dto.Items = itemsDTO

	return dto, nil
}