package repositories

import (
	"context"

	"Veritasbackend/ent"
	"Veritasbackend/ent/supplier"
)

type SupplierRepository interface {
	FindAll(ctx context.Context, tenantID int, limit, offset int) ([]*ent.Supplier, int, error)
	FindByID(ctx context.Context, id int) (*ent.Supplier, error)
	Create(ctx context.Context, tenantID int, name, email, phone, address, rucNit string) (*ent.Supplier, error)
	Update(ctx context.Context, id int, name, email, phone, address, rucNit string) (*ent.Supplier, error)
	Delete(ctx context.Context, id int) error
	CountByTenant(ctx context.Context, tenantID int) (int, error)
}

type supplierRepository struct {
	client *ent.Client
}

func NewSupplierRepository(client *ent.Client) SupplierRepository {
	return &supplierRepository{client: client}
}

func (r *supplierRepository) FindAll(ctx context.Context, tenantID int, limit, offset int) ([]*ent.Supplier, int, error) {
	query := r.client.Supplier.
		Query().
		Where(supplier.TenantIDEQ(tenantID))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	suppliers, err := query.
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(supplier.FieldCreatedAt)).
		All(ctx)

	return suppliers, total, err
}

func (r *supplierRepository) FindByID(ctx context.Context, id int) (*ent.Supplier, error) {
	return r.client.Supplier.
		Query().
		Where(supplier.IDEQ(id)).
		Only(ctx)
}

func (r *supplierRepository) Create(ctx context.Context, tenantID int, name, email, phone, address, rucNit string) (*ent.Supplier, error) {
	builder := r.client.Supplier.
		Create().
		SetTenantID(tenantID).
		SetName(name)

	if email != "" {
		builder.SetEmail(email)
	}
	if phone != "" {
		builder.SetPhone(phone)
	}
	if address != "" {
		builder.SetAddress(address)
	}
	if rucNit != "" {
		builder.SetRucNit(rucNit)
	}

	return builder.Save(ctx)
}

func (r *supplierRepository) Update(ctx context.Context, id int, name, email, phone, address, rucNit string) (*ent.Supplier, error) {
	builder := r.client.Supplier.
		UpdateOneID(id).
		SetName(name)

	if email != "" {
		builder.SetEmail(email)
	} else {
		builder.ClearEmail()
	}
	if phone != "" {
		builder.SetPhone(phone)
	} else {
		builder.ClearPhone()
	}
	if address != "" {
		builder.SetAddress(address)
	} else {
		builder.ClearAddress()
	}
	if rucNit != "" {
		builder.SetRucNit(rucNit)
	} else {
		builder.ClearRucNit()
	}

	return builder.Save(ctx)
}

func (r *supplierRepository) Delete(ctx context.Context, id int) error {
	return r.client.Supplier.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *supplierRepository) CountByTenant(ctx context.Context, tenantID int) (int, error) {
	return r.client.Supplier.
		Query().
		Where(supplier.TenantIDEQ(tenantID)).
		Count(ctx)
}