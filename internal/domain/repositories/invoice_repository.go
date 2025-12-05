package repositories

import (
	"context"
	"strconv"
	"time"

	"Veritasbackend/ent"
	"Veritasbackend/ent/invoice"
	"Veritasbackend/ent/invoiceitem"
	"Veritasbackend/ent/product"
)

type InvoiceItem struct {
	ProductID int
	Quantity  int
	UnitPrice float64
	Subtotal  float64
}

type InvoiceRepository interface {
	CountByTenant(ctx context.Context, tenantID int) (int, error)
	SumTotalByTenant(ctx context.Context, tenantID int, startDate, endDate time.Time) (float64, error)
	CountByTenantAndDateRange(ctx context.Context, tenantID int, startDate, endDate time.Time) (int, error)
	Create(ctx context.Context, tenantID, userID int, total float64, items []InvoiceItem) (*ent.Invoice, error)
	FindByID(ctx context.Context, id int) (*ent.Invoice, []*ent.InvoiceItem, error)
	FindAll(ctx context.Context, tenantID int, limit, offset int) ([]*ent.Invoice, int, error)
	SearchProducts(ctx context.Context, tenantID int, query string) ([]*ent.Product, error)
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

func (r *invoiceRepository) Create(ctx context.Context, tenantID, userID int, total float64, items []InvoiceItem) (*ent.Invoice, error) {
	// Usar transacción para crear factura e items
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	// Crear la factura
	inv, err := tx.Invoice.
		Create().
		SetTenantID(tenantID).
		SetUserID(userID).
		SetTotal(total).
		SetStatus("pending").
		Save(ctx)
	if err != nil {
		return nil, rollback(tx, err)
	}

	// Crear los items de la factura
	// Nota: Después de regenerar Ent, InvoiceItem estará disponible
	// Por ahora usamos una estructura temporal
	for _, item := range items {
		_, err = tx.InvoiceItem.
			Create().
			SetInvoiceID(inv.ID).
			SetProductID(item.ProductID).
			SetQuantity(item.Quantity).
			SetUnitPrice(item.UnitPrice).
			SetSubtotal(item.Subtotal).
			Save(ctx)
		if err != nil {
			return nil, rollback(tx, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return inv, nil
}

func (r *invoiceRepository) FindByID(ctx context.Context, id int) (*ent.Invoice, []*ent.InvoiceItem, error) {
	inv, err := r.client.Invoice.
		Query().
		Where(invoice.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	items, err := r.client.InvoiceItem.
		Query().
		Where(invoiceitem.InvoiceIDEQ(id)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	return inv, items, nil
}

func (r *invoiceRepository) FindAll(ctx context.Context, tenantID int, limit, offset int) ([]*ent.Invoice, int, error) {
	query := r.client.Invoice.
		Query().
		Where(invoice.TenantIDEQ(tenantID))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	invoices, err := query.
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(invoice.FieldCreatedAt)).
		All(ctx)

	return invoices, total, err
}

func (r *invoiceRepository) SearchProducts(ctx context.Context, tenantID int, searchQuery string) ([]*ent.Product, error) {
	// Buscar por nombre o SKU que contenga el query
	// También intentar buscar por ID si el query es numérico
	query := r.client.Product.
		Query().
		Where(product.TenantIDEQ(tenantID))

	// Si el query es numérico, intentar buscar por ID también
	if id, err := strconv.Atoi(searchQuery); err == nil {
		query = query.Where(
			product.Or(
				product.IDEQ(id),
				product.NameContainsFold(searchQuery),
				product.SkuContainsFold(searchQuery),
			),
		)
	} else {
		query = query.Where(
			product.Or(
				product.NameContainsFold(searchQuery),
				product.SkuContainsFold(searchQuery),
			),
		)
	}

	return query.Limit(20).All(ctx)
}

// rollback es una función auxiliar para hacer rollback de transacciones
func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = rerr
	}
	return err
}

