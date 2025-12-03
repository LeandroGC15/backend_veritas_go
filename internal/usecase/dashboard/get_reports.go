package dashboard

import (
	"context"
	"time"

	"Veritasbackend/internal/domain/repositories"
)

type GetReportsUseCase struct {
	invoiceRepo repositories.InvoiceRepository
}

func NewGetReportsUseCase(invoiceRepo repositories.InvoiceRepository) *GetReportsUseCase {
	return &GetReportsUseCase{
		invoiceRepo: invoiceRepo,
	}
}

type ReportRequest struct {
	Period    string `json:"period"` // daily, weekly, monthly
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type ReportData struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

type ReportResponse struct {
	Period string       `json:"period"`
	Data   []ReportData `json:"data"`
}

func (uc *GetReportsUseCase) Execute(ctx context.Context, tenantID int, req ReportRequest) (*ReportResponse, error) {
	// Parsear fechas
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		// Si no se proporciona, usar fecha por defecto
		startDate = time.Now().AddDate(0, -1, 0) // Hace un mes
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		endDate = time.Now()
	}

	// Obtener datos del repositorio
	total, err := uc.invoiceRepo.SumTotalByTenant(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	count, err := uc.invoiceRepo.CountByTenantAndDateRange(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Por ahora, retornamos un reporte simple
	// En el futuro se puede agregar agrupación por período
	data := []ReportData{
		{
			Date:  startDate.Format("2006-01-02"),
			Value: total,
			Count: count,
		},
	}

	return &ReportResponse{
		Period: req.Period,
		Data:   data,
	}, nil
}

