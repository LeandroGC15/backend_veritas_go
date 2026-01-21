package supplier

import (
	"context"

	"Veritasbackend/ent"
	"Veritasbackend/internal/domain/repositories"
)

type CreateSupplierUseCase struct {
	supplierRepo repositories.SupplierRepository
}

func NewCreateSupplierUseCase(supplierRepo repositories.SupplierRepository) *CreateSupplierUseCase {
	return &CreateSupplierUseCase{
		supplierRepo: supplierRepo,
	}
}

type CreateSupplierRequest struct {
	Name    string  `json:"name"`
	Email   *string `json:"email,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Address *string `json:"address,omitempty"`
	RucNit  *string `json:"rucNit,omitempty"`
}

type SupplierDTO struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Address   *string `json:"address,omitempty"`
	RucNit    *string `json:"rucNit,omitempty"`
	TenantID  int     `json:"tenantId"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

func convertSupplierToDTO(supplier *ent.Supplier) *SupplierDTO {
	var email, phone, address, rucNit *string

	if supplier.Email != "" {
		email = &supplier.Email
	}
	if supplier.Phone != "" {
		phone = &supplier.Phone
	}
	if supplier.Address != "" {
		address = &supplier.Address
	}
	if supplier.RucNit != "" {
		rucNit = &supplier.RucNit
	}

	return &SupplierDTO{
		ID:        supplier.ID,
		Name:      supplier.Name,
		Email:     email,
		Phone:     phone,
		Address:   address,
		RucNit:    rucNit,
		TenantID:  supplier.TenantID,
		CreatedAt: supplier.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: supplier.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (uc *CreateSupplierUseCase) Execute(ctx context.Context, tenantID int, req CreateSupplierRequest) (*SupplierDTO, error) {
	name := req.Name
	email := ""
	if req.Email != nil {
		email = *req.Email
	}
	phone := ""
	if req.Phone != nil {
		phone = *req.Phone
	}
	address := ""
	if req.Address != nil {
		address = *req.Address
	}
	rucNit := ""
	if req.RucNit != nil {
		rucNit = *req.RucNit
	}

	supplier, err := uc.supplierRepo.Create(ctx, tenantID, name, email, phone, address, rucNit)
	if err != nil {
		return nil, err
	}

	return convertSupplierToDTO(supplier), nil
}