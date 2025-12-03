package auth

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
	pkg_errors "Veritasbackend/pkg/errors"
)

type GetCurrentUserUseCase struct {
	userRepo repositories.UserRepository
}

func NewGetCurrentUserUseCase(userRepo repositories.UserRepository) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetCurrentUserUseCase) Execute(ctx context.Context, userID int) (*UserDTO, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, pkg_errors.ErrNotFound
	}

	return &UserDTO{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

