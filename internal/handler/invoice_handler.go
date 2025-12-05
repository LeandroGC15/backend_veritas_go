package handler

import (
	"net/http"
	"strconv"

	"Veritasbackend/internal/usecase/invoice"
	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	createInvoiceUseCase  *invoice.CreateInvoiceUseCase
	listInvoicesUseCase   *invoice.ListInvoicesUseCase
	getInvoiceUseCase     *invoice.GetInvoiceUseCase
	searchProductsUseCase *invoice.SearchProductsUseCase
}

func NewInvoiceHandler(
	createInvoiceUseCase *invoice.CreateInvoiceUseCase,
	listInvoicesUseCase *invoice.ListInvoicesUseCase,
	getInvoiceUseCase *invoice.GetInvoiceUseCase,
	searchProductsUseCase *invoice.SearchProductsUseCase,
) *InvoiceHandler {
	return &InvoiceHandler{
		createInvoiceUseCase:  createInvoiceUseCase,
		listInvoicesUseCase:   listInvoicesUseCase,
		getInvoiceUseCase:     getInvoiceUseCase,
		searchProductsUseCase: searchProductsUseCase,
	}
}

func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")
	userID, _ := c.Get("userID")

	var req invoice.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice, err := h.createInvoiceUseCase.Execute(c.Request.Context(), tenantID.(int), userID.(int), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"invoice": invoice})
}

func (h *InvoiceHandler) ListInvoices(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")

	req := invoice.ListInvoicesRequest{
		Page:  1,
		Limit: 20,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			req.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}

	response, err := h.listInvoicesUseCase.Execute(c.Request.Context(), tenantID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	invoice, err := h.getInvoiceUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

func (h *InvoiceHandler) SearchProducts(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusOK, gin.H{"products": []invoice.ProductDTO{}})
		return
	}

	response, err := h.searchProductsUseCase.Execute(c.Request.Context(), tenantID.(int), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

