package repositories

import (
	"context"

	"Veritasbackend/ent"
	"Veritasbackend/ent/product"
)

type ProductRepository interface {
	FindAll(ctx context.Context, tenantID int, limit, offset int) ([]*ent.Product, int, error)
	FindByID(ctx context.Context, id int) (*ent.Product, error)
	Create(ctx context.Context, tenantID int, name, description, sku string, price float64, stock int) (*ent.Product, error)
	Update(ctx context.Context, id int, name, description, sku string, price float64, stock int) (*ent.Product, error)
	Delete(ctx context.Context, id int) error
	CountByTenant(ctx context.Context, tenantID int) (int, error)
}

type productRepository struct {
	client *ent.Client
}

func NewProductRepository(client *ent.Client) ProductRepository {
	return &productRepository{client: client}
}

func (r *productRepository) FindAll(ctx context.Context, tenantID int, limit, offset int) ([]*ent.Product, int, error) {
	query := r.client.Product.
		Query().
		Where(product.TenantIDEQ(tenantID))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	products, err := query.
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(product.FieldCreatedAt)).
		All(ctx)

	return products, total, err
}

func (r *productRepository) FindByID(ctx context.Context, id int) (*ent.Product, error) {
	return r.client.Product.
		Query().
		Where(product.IDEQ(id)).
		Only(ctx)
}

func (r *productRepository) Create(ctx context.Context, tenantID int, name, description, sku string, price float64, stock int) (*ent.Product, error) {
	builder := r.client.Product.
		Create().
		SetTenantID(tenantID).
		SetName(name).
		SetPrice(price).
		SetStock(stock)

	if description != "" {
		builder.SetDescription(description)
	}
	if sku != "" {
		builder.SetSku(sku)
	}

	return builder.Save(ctx)
}

func (r *productRepository) Update(ctx context.Context, id int, name, description, sku string, price float64, stock int) (*ent.Product, error) {
	builder := r.client.Product.
		UpdateOneID(id).
		SetName(name).
		SetPrice(price).
		SetStock(stock)

	if description != "" {
		builder.SetDescription(description)
	}
	if sku != "" {
		builder.SetSku(sku)
	}

	return builder.Save(ctx)
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	return r.client.Product.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *productRepository) CountByTenant(ctx context.Context, tenantID int) (int, error) {
	return r.client.Product.
		Query().
		Where(product.TenantIDEQ(tenantID)).
		Count(ctx)
}

