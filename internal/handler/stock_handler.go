package handler

import (
	"net/http"
	"strconv"

	"Veritasbackend/internal/usecase/stock"
	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	listProductsUseCase   *stock.ListProductsUseCase
	createProductUseCase  *stock.CreateProductUseCase
	updateProductUseCase  *stock.UpdateProductUseCase
	deleteProductUseCase  *stock.DeleteProductUseCase
	uploadProductsUseCase *stock.UploadProductsUseCase
}

func NewStockHandler(
	listProductsUseCase *stock.ListProductsUseCase,
	createProductUseCase *stock.CreateProductUseCase,
	updateProductUseCase *stock.UpdateProductUseCase,
	deleteProductUseCase *stock.DeleteProductUseCase,
	uploadProductsUseCase *stock.UploadProductsUseCase,
) *StockHandler {
	return &StockHandler{
		listProductsUseCase:   listProductsUseCase,
		createProductUseCase:  createProductUseCase,
		updateProductUseCase:  updateProductUseCase,
		deleteProductUseCase:  deleteProductUseCase,
		uploadProductsUseCase: uploadProductsUseCase,
	}
}

func (h *StockHandler) ListProducts(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")

	req := stock.ListProductsRequest{
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

	response, err := h.listProductsUseCase.Execute(c.Request.Context(), tenantID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *StockHandler) CreateProduct(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")

	var req stock.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.createProductUseCase.Execute(c.Request.Context(), tenantID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": product})
}

func (h *StockHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req stock.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.updateProductUseCase.Execute(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *StockHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.deleteProductUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (h *StockHandler) UploadProducts(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	result, err := h.uploadProductsUseCase.Execute(c.Request.Context(), tenantID.(int), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

