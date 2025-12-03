package repositories

import (
	"context"

	"Veritasbackend/ent"
	"Veritasbackend/ent/user"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*ent.User, error)
	FindByID(ctx context.Context, id int) (*ent.User, error)
	Create(ctx context.Context, email, password, name, role string, tenantID int) (*ent.User, error)
}

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.
		Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*ent.User, error) {
	return r.client.User.
		Query().
		Where(user.IDEQ(id)).
		Only(ctx)
}

func (r *userRepository) Create(ctx context.Context, email, password, name, role string, tenantID int) (*ent.User, error) {
	return r.client.User.
		Create().
		SetEmail(email).
		SetPassword(password).
		SetName(name).
		SetRole(role).
		SetTenantID(tenantID).
		Save(ctx)
}

