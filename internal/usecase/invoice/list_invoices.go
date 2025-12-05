package invoice

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
)

type ListInvoicesUseCase struct {
	invoiceRepo repositories.InvoiceRepository
}

func NewListInvoicesUseCase(invoiceRepo repositories.InvoiceRepository) *ListInvoicesUseCase {
	return &ListInvoicesUseCase{
		invoiceRepo: invoiceRepo,
	}
}

type ListInvoicesRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type InvoiceSummaryDTO struct {
	ID        int     `json:"id"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
	UserID    int     `json:"userId"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type ListInvoicesResponse struct {
	Invoices []InvoiceSummaryDTO `json:"invoices"`
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	Limit    int                 `json:"limit"`
}

func (uc *ListInvoicesUseCase) Execute(ctx context.Context, tenantID int, req ListInvoicesRequest) (*ListInvoicesResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	offset := (req.Page - 1) * req.Limit

	invoices, total, err := uc.invoiceRepo.FindAll(ctx, tenantID, req.Limit, offset)
	if err != nil {
		return nil, err
	}

	invoiceDTOs := make([]InvoiceSummaryDTO, len(invoices))
	for i, inv := range invoices {
		invoiceDTOs[i] = InvoiceSummaryDTO{
			ID:        inv.ID,
			Total:     inv.Total,
			Status:    inv.Status,
			UserID:    inv.UserID,
			CreatedAt: inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: inv.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &ListInvoicesResponse{
		Invoices: invoiceDTOs,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

