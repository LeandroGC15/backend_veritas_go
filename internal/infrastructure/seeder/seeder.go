package seeder

import (
	"context"
	"log"

	"Veritasbackend/ent"
	"Veritasbackend/ent/tenant"
	"Veritasbackend/ent/user"
	"golang.org/x/crypto/bcrypt"
)

type Seeder struct {
	client *ent.Client
}

func NewSeeder(client *ent.Client) *Seeder {
	return &Seeder{client: client}
}

func (s *Seeder) Seed(ctx context.Context) error {
	log.Println("🌱 Starting database seeding...")

	// 1. Crear tenants
	tenants, err := s.seedTenants(ctx)
	if err != nil {
		return err
	}

	// 2. Crear usuarios
	if err := s.seedUsers(ctx, tenants); err != nil {
		return err
	}

	log.Println("✅ Seeding completed successfully!")
	return nil
}

func (s *Seeder) seedTenants(ctx context.Context) (map[string]*ent.Tenant, error) {
	log.Println("📦 Seeding tenants...")

	tenants := []struct {
		name   string
		slug   string
		domain string
	}{
		{
			name:   "Empresa Demo",
			slug:   "demo",
			domain: "demo.veritas.com",
		},
		{
			name:   "Acme Corporation",
			slug:   "acme",
			domain: "acme.com",
		},
		{
			name:   "Tech Solutions",
			slug:   "tech",
			domain: "techsolutions.com",
		},
	}

	tenantMap := make(map[string]*ent.Tenant)

	for _, t := range tenants {
		// Verificar si el tenant ya existe
		existing, err := s.client.Tenant.Query().
			Where(tenant.SlugEQ(t.slug)).
			Only(ctx)

		if err != nil && !ent.IsNotFound(err) {
			return nil, err
		}

		if existing != nil {
			log.Printf("  ⏭️  Tenant '%s' already exists, skipping...", t.name)
			tenantMap[t.slug] = existing
			continue
		}

		// Crear nuevo tenant
		newTenant, err := s.client.Tenant.
			Create().
			SetName(t.name).
			SetSlug(t.slug).
			SetDomain(t.domain).
			Save(ctx)

		if err != nil {
			return nil, err
		}

		log.Printf("  ✅ Created tenant: %s (slug: %s)", t.name, t.slug)
		tenantMap[t.slug] = newTenant
	}

	return tenantMap, nil
}

func (s *Seeder) seedUsers(ctx context.Context, tenants map[string]*ent.Tenant) error {
	log.Println("👥 Seeding users...")

	users := []struct {
		email    string
		password string
		name     string
		role     string
		tenant   string
	}{
		// Usuarios para tenant "demo"
		{
			email:    "admin@demo.veritas.com",
			password: "admin123",
			name:     "Administrador Demo",
			role:     "admin",
			tenant:   "demo",
		},
		{
			email:    "user@demo.veritas.com",
			password: "user123",
			name:     "Usuario Demo",
			role:     "user",
			tenant:   "demo",
		},
		// Usuarios para tenant "acme"
		{
			email:    "admin@acme.com",
			password: "admin123",
			name:     "Admin Acme",
			role:     "admin",
			tenant:   "acme",
		},
		{
			email:    "manager@acme.com",
			password: "manager123",
			name:     "Manager Acme",
			role:     "manager",
			tenant:   "acme",
		},
		// Usuarios para tenant "tech"
		{
			email:    "admin@techsolutions.com",
			password: "admin123",
			name:     "Admin Tech",
			role:     "admin",
			tenant:   "tech",
		},
		{
			email:    "user@techsolutions.com",
			password: "user123",
			name:     "Usuario Tech",
			role:     "user",
			tenant:   "tech",
		},
	}

	for _, u := range users {
		tenant, ok := tenants[u.tenant]
		if !ok {
			log.Printf("  ⚠️  Tenant '%s' not found for user %s, skipping...", u.tenant, u.email)
			continue
		}

		// Verificar si el usuario ya existe
		existing, err := s.client.User.Query().
			Where(user.EmailEQ(u.email)).
			Only(ctx)

		if err != nil && !ent.IsNotFound(err) {
			return err
		}

		if existing != nil {
			log.Printf("  ⏭️  User '%s' already exists, skipping...", u.email)
			continue
		}

		// Hashear password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Crear usuario
		_, err = s.client.User.
			Create().
			SetEmail(u.email).
			SetPassword(string(hashedPassword)).
			SetName(u.name).
			SetRole(u.role).
			SetTenantID(tenant.ID).
			Save(ctx)

		if err != nil {
			return err
		}

		log.Printf("  ✅ Created user: %s (role: %s, tenant: %s)", u.email, u.role, u.tenant)
	}

	return nil
}

