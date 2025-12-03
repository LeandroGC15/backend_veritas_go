package handler

import (
	"net/http"

	"Veritasbackend/internal/usecase/dashboard"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	getMetricsUseCase *dashboard.GetMetricsUseCase
	getReportsUseCase *dashboard.GetReportsUseCase
}

func NewDashboardHandler(getMetricsUseCase *dashboard.GetMetricsUseCase, getReportsUseCase *dashboard.GetReportsUseCase) *DashboardHandler {
	return &DashboardHandler{
		getMetricsUseCase: getMetricsUseCase,
		getReportsUseCase: getReportsUseCase,
	}
}

func (h *DashboardHandler) GetMetrics(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")

	metrics, err := h.getMetricsUseCase.Execute(c.Request.Context(), tenantID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *DashboardHandler) GetReports(c *gin.Context) {
	tenantID, _ := c.Get("tenantID")

	req := dashboard.ReportRequest{
		Period:    c.Query("period"),
		StartDate: c.Query("startDate"),
		EndDate:   c.Query("endDate"),
	}

	reports, err := h.getReportsUseCase.Execute(c.Request.Context(), tenantID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

