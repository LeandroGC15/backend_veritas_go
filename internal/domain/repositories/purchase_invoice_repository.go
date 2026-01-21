package repositories

import (
	"context"
	"time"

	"Veritasbackend/ent"
	"Veritasbackend/ent/purchaseinvoice"
	"Veritasbackend/ent/purchaseinvoiceitem"
)

type PurchaseInvoiceRepository interface {
	FindAll(ctx context.Context, tenantID int, limit, offset int) ([]*ent.PurchaseInvoice, int, error)
	FindByID(ctx context.Context, id int) (*ent.PurchaseInvoice, error)
	Create(ctx context.Context, tenantID, supplierID, userID int, invoiceNumber string, total float64, paymentMethod *string, dueDate *time.Time) (*ent.PurchaseInvoice, error)
	Update(ctx context.Context, id int, status string, paidAmount float64) (*ent.PurchaseInvoice, error)
	Delete(ctx context.Context, id int) error
	FindBySupplierID(ctx context.Context, supplierID int) ([]*ent.PurchaseInvoice, error)
}

type PurchaseInvoiceItemRepository interface {
	FindByPurchaseInvoiceID(ctx context.Context, purchaseInvoiceID int) ([]*ent.PurchaseInvoiceItem, error)
	Create(ctx context.Context, purchaseInvoiceID, productID, quantity int, unitCost, subtotal float64) (*ent.PurchaseInvoiceItem, error)
	CreateBulk(ctx context.Context, items []*ent.PurchaseInvoiceItem) error
	DeleteByPurchaseInvoiceID(ctx context.Context, purchaseInvoiceID int) error
}

type purchaseInvoiceRepository struct {
	client *ent.Client
}

func NewPurchaseInvoiceRepository(client *ent.Client) PurchaseInvoiceRepository {
	return &purchaseInvoiceRepository{client: client}
}

func (r *purchaseInvoiceRepository) FindAll(ctx context.Context, tenantID int, limit, offset int) ([]*ent.PurchaseInvoice, int, error) {
	query := r.client.PurchaseInvoice.
		Query().
		Where(purchaseinvoice.TenantIDEQ(tenantID))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	invoices, err := query.
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(purchaseinvoice.FieldCreatedAt)).
		All(ctx)

	return invoices, total, err
}

func (r *purchaseInvoiceRepository) FindByID(ctx context.Context, id int) (*ent.PurchaseInvoice, error) {
	return r.client.PurchaseInvoice.
		Query().
		Where(purchaseinvoice.IDEQ(id)).
		Only(ctx)
}

func (r *purchaseInvoiceRepository) Create(ctx context.Context, tenantID, supplierID, userID int, invoiceNumber string, total float64, paymentMethod *string, dueDate *time.Time) (*ent.PurchaseInvoice, error) {
	builder := r.client.PurchaseInvoice.
		Create().
		SetTenantID(tenantID).
		SetSupplierID(supplierID).
		SetUserID(userID).
		SetInvoiceNumber(invoiceNumber).
		SetTotal(total).
		SetStatus("pending").
		SetPaidAmount(0)

	if paymentMethod != nil {
		builder.SetPaymentMethod(*paymentMethod)
	}

	if dueDate != nil {
		builder.SetDueDate(*dueDate)
	}

	return builder.Save(ctx)
}

func (r *purchaseInvoiceRepository) Update(ctx context.Context, id int, status string, paidAmount float64) (*ent.PurchaseInvoice, error) {
	return r.client.PurchaseInvoice.
		UpdateOneID(id).
		SetStatus(status).
		SetPaidAmount(paidAmount).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

func (r *purchaseInvoiceRepository) Delete(ctx context.Context, id int) error {
	return r.client.PurchaseInvoice.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *purchaseInvoiceRepository) FindBySupplierID(ctx context.Context, supplierID int) ([]*ent.PurchaseInvoice, error) {
	return r.client.PurchaseInvoice.
		Query().
		Where(purchaseinvoice.SupplierIDEQ(supplierID)).
		Order(ent.Desc(purchaseinvoice.FieldCreatedAt)).
		All(ctx)
}

type purchaseInvoiceItemRepository struct {
	client *ent.Client
}

func NewPurchaseInvoiceItemRepository(client *ent.Client) PurchaseInvoiceItemRepository {
	return &purchaseInvoiceItemRepository{client: client}
}

func (r *purchaseInvoiceItemRepository) FindByPurchaseInvoiceID(ctx context.Context, purchaseInvoiceID int) ([]*ent.PurchaseInvoiceItem, error) {
	return r.client.PurchaseInvoiceItem.
		Query().
		Where(purchaseinvoiceitem.PurchaseInvoiceIDEQ(purchaseInvoiceID)).
		All(ctx)
}

func (r *purchaseInvoiceItemRepository) Create(ctx context.Context, purchaseInvoiceID, productID, quantity int, unitCost, subtotal float64) (*ent.PurchaseInvoiceItem, error) {
	return r.client.PurchaseInvoiceItem.
		Create().
		SetPurchaseInvoiceID(purchaseInvoiceID).
		SetProductID(productID).
		SetQuantity(quantity).
		SetUnitCost(unitCost).
		SetSubtotal(subtotal).
		Save(ctx)
}

func (r *purchaseInvoiceItemRepository) CreateBulk(ctx context.Context, items []*ent.PurchaseInvoiceItem) error {
	// For bulk creation, we'll create them one by one for now
	// In a real implementation, you might want to use a transaction
	for _, item := range items {
		_, err := r.client.PurchaseInvoiceItem.
			Create().
			SetPurchaseInvoiceID(item.PurchaseInvoiceID).
			SetProductID(item.ProductID).
			SetQuantity(item.Quantity).
			SetUnitCost(item.UnitCost).
			SetSubtotal(item.Subtotal).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *purchaseInvoiceItemRepository) DeleteByPurchaseInvoiceID(ctx context.Context, purchaseInvoiceID int) error {
	_, err := r.client.PurchaseInvoiceItem.
		Delete().
		Where(purchaseinvoiceitem.PurchaseInvoiceIDEQ(purchaseInvoiceID)).
		Exec(ctx)
	return err
}