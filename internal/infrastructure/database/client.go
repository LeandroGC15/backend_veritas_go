package database

import (
	"context"
	"fmt"
	"log"

	"Veritasbackend/ent"
	"Veritasbackend/internal/infrastructure/config"

	_ "github.com/lib/pq"
)

func NewClient(cfg *config.Config) (*ent.Client, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	// Ejecutar migraciones automáticamente
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Printf("Warning: failed creating schema resources: %v", err)
	}

	return client, nil
}
