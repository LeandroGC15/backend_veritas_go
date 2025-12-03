package main

import (
	"context"
	"log"

	"Veritasbackend/internal/infrastructure/config"
	"Veritasbackend/internal/infrastructure/database"
	"Veritasbackend/internal/infrastructure/seeder"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	// Conectar a la base de datos
	client, err := database.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Ejecutar seeder
	seed := seeder.NewSeeder(client)
	if err := seed.Seed(ctx); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("✅ Database seeded successfully!")
}
