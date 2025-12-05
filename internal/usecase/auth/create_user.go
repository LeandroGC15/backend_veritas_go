package auth

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"Veritasbackend/internal/domain/repositories"
	pkg_errors "Veritasbackend/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase struct {
	userRepo   repositories.UserRepository
	tenantRepo repositories.TenantRepository
}

func NewCreateUserUseCase(userRepo repositories.UserRepository, tenantRepo repositories.TenantRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// generateUniqueSlug genera un slug único basado en el email
func generateUniqueSlug(email string) string {
	// Extraer la parte antes y después del @
	parts := strings.Split(email, "@")
	username := strings.ToLower(parts[0])
	domain := strings.ToLower(parts[1])
	
	// Limpiar el username: remover caracteres especiales y espacios
	username = strings.ReplaceAll(username, ".", "-")
	username = strings.ReplaceAll(username, "_", "-")
	username = strings.ReplaceAll(username, " ", "-")
	
	// Limpiar el dominio: remover puntos
	domain = strings.ReplaceAll(domain, ".", "-")
	
	// Agregar timestamp para garantizar unicidad
	now := time.Now()
	timestamp := now.Format("20060102-150405")
	
	return username + "-" + domain + "-" + timestamp
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	log.Printf("📝 CreateUserUseCase: Iniciando creación de usuario - Email: %s, Name: %s, Role: %s", req.Email, req.Name, req.Role)

	// Verificar si el usuario ya existe
	log.Printf("🔍 CreateUserUseCase: Verificando si el usuario ya existe - Email: %s", req.Email)
	existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		log.Printf("❌ CreateUserUseCase: Usuario ya existe - Email: %s", req.Email)
		return nil, errors.New("user already exists")
	}
	log.Printf("✅ CreateUserUseCase: Email disponible")

	// Validar rol
	log.Printf("🔍 CreateUserUseCase: Validando rol: %s", req.Role)
	validRoles := map[string]bool{"admin": true, "manager": true, "user": true}
	if !validRoles[req.Role] {
		log.Printf("❌ CreateUserUseCase: Rol inválido: %s", req.Role)
		return nil, pkg_errors.ErrInvalidInput
	}
	log.Printf("✅ CreateUserUseCase: Rol válido")

	// Crear un nuevo tenant único para este usuario
	log.Printf("🏢 CreateUserUseCase: Creando nuevo tenant para el usuario...")
	tenantName := req.Name + " - " + req.Email
	tenantSlug := generateUniqueSlug(req.Email)
	tenantDomain := strings.Split(req.Email, "@")[1]
	
	// Verificar que el slug no exista (aunque es muy improbable con timestamp)
	existingTenant, _ := uc.tenantRepo.FindBySlug(ctx, tenantSlug)
	if existingTenant != nil {
		// Si existe, agregar más aleatoriedad con nanosegundos
		now := time.Now()
		nanos := now.Format("150405.000000000")
		nanos = strings.ReplaceAll(nanos, ".", "-")
		tenantSlug = tenantSlug + "-" + nanos
	}
	
	newTenant, err := uc.tenantRepo.Create(ctx, tenantName, tenantSlug, tenantDomain)
	if err != nil {
		log.Printf("❌ CreateUserUseCase: Error al crear tenant: %v", err)
		return nil, errors.New("failed to create tenant")
	}
	log.Printf("✅ CreateUserUseCase: Tenant creado exitosamente - ID: %d, Slug: %s", newTenant.ID, tenantSlug)

	// Hashear password
	log.Printf("🔐 CreateUserUseCase: Hasheando contraseña...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ CreateUserUseCase: Error al hashear contraseña: %v", err)
		return nil, errors.New("failed to hash password")
	}
	log.Printf("✅ CreateUserUseCase: Contraseña hasheada correctamente")

	// Crear usuario con el nuevo tenant
	log.Printf("💾 CreateUserUseCase: Creando usuario en la base de datos con tenant ID: %d...", newTenant.ID)
	user, err := uc.userRepo.Create(ctx, req.Email, string(hashedPassword), req.Name, req.Role, newTenant.ID)
	if err != nil {
		log.Printf("❌ CreateUserUseCase: Error al crear usuario en BD: %v", err)
		return nil, errors.New("failed to create user")
	}

	log.Printf("✅ CreateUserUseCase: Usuario creado exitosamente - ID: %d, Email: %s, TenantID: %d", user.ID, user.Email, newTenant.ID)

	return &UserDTO{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

