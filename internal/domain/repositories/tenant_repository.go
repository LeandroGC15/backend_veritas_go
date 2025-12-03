package repositories

import (
	"context"

	"Veritasbackend/ent"
	"Veritasbackend/ent/tenant"
)

type TenantRepository interface {
	FindByID(ctx context.Context, id int) (*ent.Tenant, error)
	FindBySlug(ctx context.Context, slug string) (*ent.Tenant, error)
	Create(ctx context.Context, name, slug, domain string) (*ent.Tenant, error)
}

type tenantRepository struct {
	client *ent.Client
}

func NewTenantRepository(client *ent.Client) TenantRepository {
	return &tenantRepository{client: client}
}

func (r *tenantRepository) FindByID(ctx context.Context, id int) (*ent.Tenant, error) {
	return r.client.Tenant.
		Query().
		Where(tenant.IDEQ(id)).
		Only(ctx)
}

func (r *tenantRepository) FindBySlug(ctx context.Context, slug string) (*ent.Tenant, error) {
	return r.client.Tenant.
		Query().
		Where(tenant.SlugEQ(slug)).
		Only(ctx)
}

func (r *tenantRepository) Create(ctx context.Context, name, slug, domain string) (*ent.Tenant, error) {
	return r.client.Tenant.
		Create().
		SetName(name).
		SetSlug(slug).
		SetDomain(domain).
		Save(ctx)
}

