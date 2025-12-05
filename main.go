package main

import (
	"log"

	"Veritasbackend/internal/domain/repositories"
	"Veritasbackend/internal/handler"
	"Veritasbackend/internal/infrastructure/config"
	"Veritasbackend/internal/infrastructure/database"
	"Veritasbackend/internal/infrastructure/middleware"
	"Veritasbackend/internal/usecase/auth"
	"Veritasbackend/internal/usecase/dashboard"
	"Veritasbackend/internal/usecase/invoice"
	"Veritasbackend/internal/usecase/stock"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	// Conectar a la base de datos
	dbClient, err := database.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	// Inicializar repositorios
	userRepo := repositories.NewUserRepository(dbClient)
	tenantRepo := repositories.NewTenantRepository(dbClient)
	productRepo := repositories.NewProductRepository(dbClient)
	invoiceRepo := repositories.NewInvoiceRepository(dbClient)

	// Inicializar casos de uso
	loginUseCase := auth.NewLoginUseCase(userRepo, tenantRepo)
	getCurrentUserUseCase := auth.NewGetCurrentUserUseCase(userRepo)
	createUserUseCase := auth.NewCreateUserUseCase(userRepo, tenantRepo)
	getMetricsUseCase := dashboard.NewGetMetricsUseCase(productRepo, invoiceRepo)
	getReportsUseCase := dashboard.NewGetReportsUseCase(invoiceRepo)
	listProductsUseCase := stock.NewListProductsUseCase(productRepo)
	createProductUseCase := stock.NewCreateProductUseCase(productRepo)
	updateProductUseCase := stock.NewUpdateProductUseCase(productRepo)
	deleteProductUseCase := stock.NewDeleteProductUseCase(productRepo)
	uploadProductsUseCase := stock.NewUploadProductsUseCase(productRepo)
	createInvoiceUseCase := invoice.NewCreateInvoiceUseCase(invoiceRepo, productRepo)
	listInvoicesUseCase := invoice.NewListInvoicesUseCase(invoiceRepo)
	getInvoiceUseCase := invoice.NewGetInvoiceUseCase(invoiceRepo, productRepo)
	searchProductsUseCase := invoice.NewSearchProductsUseCase(invoiceRepo)

	// Inicializar handlers
	authHandler := handler.NewAuthHandler(loginUseCase, getCurrentUserUseCase, createUserUseCase)
	dashboardHandler := handler.NewDashboardHandler(getMetricsUseCase, getReportsUseCase)
	stockHandler := handler.NewStockHandler(
		listProductsUseCase,
		createProductUseCase,
		updateProductUseCase,
		deleteProductUseCase,
		uploadProductsUseCase,
	)
	invoiceHandler := handler.NewInvoiceHandler(
		createInvoiceUseCase,
		listInvoicesUseCase,
		getInvoiceUseCase,
		searchProductsUseCase,
	)

	// Configurar Gin
	if cfg.Server.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Configurar CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CORS.AllowedOrigins}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Tenant-ID"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// Rutas públicas
	api := r.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
	}

	// Rutas protegidas
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.TenantMiddleware())
	{
		// Auth
		protected.GET("/auth/me", authHandler.GetCurrentUser)

		// Admin routes - solo para usuarios con rol admin
		log.Println("🔧 Registrando ruta admin POST /api/users")
		admin := protected.Group("")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/users", authHandler.CreateUser)
			log.Println("✅ Ruta admin POST /api/users registrada correctamente")
		}

		// Dashboard
		protected.GET("/dashboard/metrics", dashboardHandler.GetMetrics)
		protected.GET("/dashboard/reports", dashboardHandler.GetReports)

		// Stock
		protected.GET("/stock", stockHandler.ListProducts)
		protected.POST("/stock", stockHandler.CreateProduct)
		protected.PUT("/stock/:id", stockHandler.UpdateProduct)
		protected.DELETE("/stock/:id", stockHandler.DeleteProduct)
		protected.POST("/stock/upload", stockHandler.UploadProducts)

		// Invoices
		protected.POST("/invoices", invoiceHandler.CreateInvoice)
		protected.GET("/invoices", invoiceHandler.ListInvoices)
		protected.GET("/invoices/:id", invoiceHandler.GetInvoice)
		protected.GET("/invoices/products/search", invoiceHandler.SearchProducts)
	}

	// Log de todas las rutas registradas para debugging
	log.Println("📋 Rutas registradas:")
	log.Println("  - POST /api/auth/login (pública)")
	log.Println("  - GET /api/auth/me (protegida)")
	log.Println("  - POST /api/users (admin)")
	log.Println("  - GET /api/dashboard/metrics (protegida)")
	log.Println("  - GET /api/dashboard/reports (protegida)")
	log.Println("  - GET /api/stock (protegida)")
	log.Println("  - POST /api/stock (protegida)")
	log.Println("  - POST /api/invoices (protegida)")
	log.Println("  - GET /api/invoices (protegida)")
	log.Println("  - GET /api/invoices/:id (protegida)")
	log.Println("  - GET /api/invoices/products/search (protegida)")
	
	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

