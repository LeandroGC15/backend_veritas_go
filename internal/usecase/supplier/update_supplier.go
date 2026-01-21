package supplier

import (
	"context"

	"Veritasbackend/internal/domain/repositories"
)

type UpdateSupplierUseCase struct {
	supplierRepo repositories.SupplierRepository
}

func NewUpdateSupplierUseCase(supplierRepo repositories.SupplierRepository) *UpdateSupplierUseCase {
	return &UpdateSupplierUseCase{
		supplierRepo: supplierRepo,
	}
}

type UpdateSupplierRequest struct {
	Name    *string `json:"name,omitempty"`
	Email   *string `json:"email,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Address *string `json:"address,omitempty"`
	RucNit  *string `json:"rucNit,omitempty"`
}

func (uc *UpdateSupplierUseCase) Execute(ctx context.Context, tenantID, supplierID int, req UpdateSupplierRequest) (*SupplierDTO, error) {
	// First verify the supplier exists and belongs to the tenant
	existingSupplier, err := uc.supplierRepo.FindByID(ctx, supplierID)
	if err != nil {
		return nil, err
	}
	if existingSupplier == nil {
		return nil, &NotFoundError{Resource: "supplier", ID: supplierID}
	}
	if existingSupplier.TenantID != tenantID {
		return nil, &ForbiddenError{Message: "Supplier does not belong to this tenant"}
	}

	// Prepare update data - use existing values if not provided
	name := existingSupplier.Name
	if req.Name != nil && *req.Name != "" {
		name = *req.Name
	}

	email := ""
	if existingSupplier.Email != "" {
		email = existingSupplier.Email
	}
	if req.Email != nil {
		if *req.Email == "" {
			email = "" // Allow clearing email
		} else {
			email = *req.Email
		}
	}

	phone := ""
	if existingSupplier.Phone != "" {
		phone = existingSupplier.Phone
	}
	if req.Phone != nil {
		if *req.Phone == "" {
			phone = "" // Allow clearing phone
		} else {
			phone = *req.Phone
		}
	}

	address := ""
	if existingSupplier.Address != "" {
		address = existingSupplier.Address
	}
	if req.Address != nil {
		if *req.Address == "" {
			address = "" // Allow clearing address
		} else {
			address = *req.Address
		}
	}

	rucNit := ""
	if existingSupplier.RucNit != "" {
		rucNit = existingSupplier.RucNit
	}
	if req.RucNit != nil {
		if *req.RucNit == "" {
			rucNit = "" // Allow clearing rucNit
		} else {
			rucNit = *req.RucNit
		}
	}

	// Update the supplier
	updatedSupplier, err := uc.supplierRepo.Update(ctx, supplierID, name, email, phone, address, rucNit)
	if err != nil {
		return nil, err
	}

	return convertSupplierToDTO(updatedSupplier), nil
}

// Error types
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return "not found"
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}