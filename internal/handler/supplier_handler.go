package handler

import (
	"net/http"
	"strconv"

	"Veritasbackend/internal/usecase/supplier"
	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	createSupplierUseCase *supplier.CreateSupplierUseCase
	listSuppliersUseCase  *supplier.ListSuppliersUseCase
	updateSupplierUseCase *supplier.UpdateSupplierUseCase
}

func NewSupplierHandler(
	createSupplierUseCase *supplier.CreateSupplierUseCase,
	listSuppliersUseCase *supplier.ListSuppliersUseCase,
	updateSupplierUseCase *supplier.UpdateSupplierUseCase,
) *SupplierHandler {
	return &SupplierHandler{
		createSupplierUseCase: createSupplierUseCase,
		listSuppliersUseCase:  listSuppliersUseCase,
		updateSupplierUseCase: updateSupplierUseCase,
	}
}

func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req supplier.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, exists := c.Get("tenantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	_ = userID // Not used in this handler, but available for logging

	result, err := h.createSupplierUseCase.Execute(c.Request.Context(), tenantID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	tenantID, exists := c.Get("tenantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	req := supplier.ListSuppliersRequest{
		Page:  page,
		Limit: limit,
	}

	result, err := h.listSuppliersUseCase.Execute(c.Request.Context(), tenantID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	var req supplier.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID, exists := c.Get("tenantID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant ID not found"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	_ = userID // Not used in this handler, but available for logging

	result, err := h.updateSupplierUseCase.Execute(c.Request.Context(), tenantID.(int), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}