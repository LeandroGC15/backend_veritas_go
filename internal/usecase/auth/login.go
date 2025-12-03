package auth

import (
	"context"
	"errors"

	"Veritasbackend/internal/domain/repositories"
	pkg_errors "Veritasbackend/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	userRepo   repositories.UserRepository
	tenantRepo repositories.TenantRepository
}

func NewLoginUseCase(userRepo repositories.UserRepository, tenantRepo repositories.TenantRepository) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	User     UserDTO     `json:"user"`
	TenantID int         `json:"tenantId"`
}

type UserDTO struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Validar email
	if req.Email == "" || req.Password == "" {
		return nil, pkg_errors.ErrInvalidInput
	}

	// Buscar usuario
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, pkg_errors.ErrUnauthorized
	}

	// Verificar password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, pkg_errors.ErrUnauthorized
	}

	// Obtener tenant
	tenant, err := uc.tenantRepo.FindByID(ctx, user.TenantID)
	if err != nil {
		return nil, errors.New("tenant not found")
	}

	return &LoginResponse{
		User: UserDTO{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
		TenantID: tenant.ID,
	}, nil
}

