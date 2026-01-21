package handler

import (
	"log"
	"net/http"
	"time"

	"Veritasbackend/internal/usecase/purchase"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	createPurchaseUseCase *purchase.CreatePurchaseUseCase
}

func NewPurchaseHandler(createPurchaseUseCase *purchase.CreatePurchaseUseCase) *PurchaseHandler {
	return &PurchaseHandler{
		createPurchaseUseCase: createPurchaseUseCase,
	}
}

type CreatePurchaseRequest struct {
	SupplierID    int                             `json:"supplierId" binding:"required"`
	InvoiceNumber string                          `json:"invoiceNumber" binding:"required"`
	PaymentMethod *string                         `json:"paymentMethod,omitempty"`
	DueDate       *string                         `json:"dueDate,omitempty"`
	Items         []CreatePurchaseItemRequest     `json:"items" binding:"required,min=1"`
}

type CreatePurchaseItemRequest struct {
	ProductID    int     `json:"productId" binding:"required"`
	Quantity     int     `json:"quantity" binding:"required,min=1"`
	UnitCost     float64 `json:"unitCost" binding:"required,min=0"`
	ProductName  string  `json:"productName"`
	ProductSku   string  `json:"productSku"`
	ProductPrice float64 `json:"productPrice"`
}

func (h *PurchaseHandler) CreatePurchase(c *gin.Context) {
	log.Printf("CreatePurchase handler called")

	tenantID, exists := c.Get("tenantID")
	if !exists {
		log.Printf("Tenant ID not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID not found"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Printf("User ID not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	log.Printf("Processing purchase for tenant %v, user %v", tenantID, userID)

	var req CreatePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	log.Printf("Request parsed successfully: supplierId=%d, invoiceNumber=%s, items=%d", req.SupplierID, req.InvoiceNumber, len(req.Items))
	for i, item := range req.Items {
		log.Printf("Parsed item %d: ProductID=%d, Quantity=%d, UnitCost=%f, ProductName='%s'", i, item.ProductID, item.Quantity, item.UnitCost, item.ProductName)
	}

	// Parse due date if provided
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format. Use YYYY-MM-DD"})
			return
		}
		dueDate = &parsed
	}

	// Convert request to use case format
	items := make([]purchase.PurchaseItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = purchase.PurchaseItemRequest{
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			UnitCost:     item.UnitCost,
			ProductName:  item.ProductName,
			ProductSku:   item.ProductSku,
			ProductPrice: item.ProductPrice,
		}
	}

	useCaseReq := purchase.CreatePurchaseRequest{
		SupplierID:    req.SupplierID,
		InvoiceNumber: req.InvoiceNumber,
		PaymentMethod: req.PaymentMethod,
		DueDate:       dueDate,
		Items:         items,
	}

	result, err := h.createPurchaseUseCase.Execute(c.Request.Context(), tenantID.(int), userID.(int), useCaseReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create purchase"})
		return
	}

	c.JSON(http.StatusCreated, result)
}