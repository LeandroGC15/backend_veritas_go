package repositories

import (
	"context"
	"time"

	"Veritasbackend/ent"
	"Veritasbackend/ent/invoice"
)

type InvoiceRepository interface {
	CountByTenant(ctx context.Context, tenantID int) (int, error)
	SumTotalByTenant(ctx context.Context, tenantID int, startDate, endDate time.Time) (float64, error)
	CountByTenantAndDateRange(ctx context.Context, tenantID int, startDate, endDate time.Time) (int, error)
}

type invoiceRepository struct {
	client *ent.Client
}

func NewInvoiceRepository(client *ent.Client) InvoiceRepository {
	return &invoiceRepository{client: client}
}

func (r *invoiceRepository) CountByTenant(ctx context.Context, tenantID int) (int, error) {
	return r.client.Invoice.
		Query().
		Where(invoice.TenantIDEQ(tenantID)).
		Count(ctx)
}

func (r *invoiceRepository) SumTotalByTenant(ctx context.Context, tenantID int, startDate, endDate time.Time) (float64, error) {
	invoices, err := r.client.Invoice.
		Query().
		Where(
			invoice.TenantIDEQ(tenantID),
			invoice.CreatedAtGTE(startDate),
			invoice.CreatedAtLTE(endDate),
		).
		All(ctx)

	if err != nil {
		return 0, err
	}

	var total float64
	for _, inv := range invoices {
		total += inv.Total
	}

	return total, nil
}

func (r *invoiceRepository) CountByTenantAndDateRange(ctx context.Context, tenantID int, startDate, endDate time.Time) (int, error) {
	return r.client.Invoice.
		Query().
		Where(
			invoice.TenantIDEQ(tenantID),
			invoice.CreatedAtGTE(startDate),
			invoice.CreatedAtLTE(endDate),
		).
		Count(ctx)
}

