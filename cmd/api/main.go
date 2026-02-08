package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kasir-api/internal/config"
	"kasir-api/internal/database"
	"kasir-api/internal/handler"
	"kasir-api/internal/repository"
	"kasir-api/internal/repository/memory"
	"kasir-api/internal/repository/postgres"
	"kasir-api/internal/service"
	"kasir-api/pkg/logger"
)

func main() {
	// Initialize logger
	logLevel := logger.NewDefault()
	logger.InitGlobalLogger(logLevel)

	// Check for help flags
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help") {
		printHelp()
		return
	}

	// Check for subcommands
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "migrate-reset" {
		runMigrateReset()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		runSeeds()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "rls" {
		runRLS()
		return
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize repositories
	var productRepo repository.ProductReader
	var productWriter repository.ProductWriter
	var categoryRepo repository.CategoryReader
	var categoryWriter repository.CategoryWriter
	var transactionWriter repository.TransactionWriter
	var reportReader repository.ReportReader
	var db *database.DB

	// Check if database is configured
	if cfg.Database.Host != "" && cfg.Database.DBName != "" {
		// Use PostgreSQL
		var err error
		db, err = database.NewPool(cfg.Database)
		if err != nil {
			logger.Error("Failed to connect to database", "error", err)
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		logger.Info("Connected to PostgreSQL database")

		// Create PostgreSQL repositories
		pgProductRepo := postgres.NewProductRepository(db.DB)
		productRepo = pgProductRepo
		productWriter = pgProductRepo

		pgCategoryRepo := postgres.NewCategoryRepository(db.DB)
		categoryRepo = pgCategoryRepo
		categoryWriter = pgCategoryRepo

		pgTransactionRepo := postgres.NewTransactionRepository(db.DB)
		transactionWriter = pgTransactionRepo

		pgReportRepo := postgres.NewReportRepository(db.DB)
		reportReader = pgReportRepo
	} else {
		// Use in-memory repositories
		logger.Info("Using in-memory storage")
		memProductRepo := memory.NewProductRepository()
		memCategoryRepo := memory.NewCategoryRepository()

		// Wire category repo to product repo for JOIN simulation
		memProductRepo.SetCategoryRepo(memCategoryRepo)

		productRepo = memProductRepo
		productWriter = memProductRepo

		categoryRepo = memCategoryRepo
		categoryWriter = memCategoryRepo
	}

	// Initialize services
	productService := service.NewProductService(productRepo, productWriter)
	categoryService := service.NewCategoryService(categoryRepo, categoryWriter)

	var transactionService *service.TransactionService
	if transactionWriter != nil {
		transactionService = service.NewTransactionService(transactionWriter)
	}

	var reportService *service.ReportService
	if reportReader != nil {
		reportService = service.NewReportService(reportReader)
	}

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	var transactionHandler *handler.TransactionHandler
	if transactionService != nil {
		transactionHandler = handler.NewTransactionHandler(transactionService)
	}

	var reportHandler *handler.ReportHandler
	if reportService != nil {
		reportHandler = handler.NewReportHandler(reportService)
	}

	var healthHandler *handler.HealthHandler
	if db != nil {
		healthHandler = handler.NewHealthHandler(db.DB)
	} else {
		healthHandler = handler.NewHealthHandler(nil)
	}

	// Setup routes
	mux := http.NewServeMux()
	handlerWithMiddleware := handler.SetupRoutes(mux, productHandler, categoryHandler, transactionHandler, reportHandler, healthHandler)

	// Create server
	server := &http.Server{
		Addr:         cfg.Server.Host + cfg.Server.Port,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		logger.Info("Starting server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Error("Server failed", "error", err)
		log.Fatalf("Server failed: %v", err)
	case sig := <-quit:
		logger.Info("Received shutdown signal", "signal", sig.String())
		logger.Info("Shutting down server...")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	} else {
		logger.Info("Server stopped gracefully")
	}
}

func printHelp() {
	fmt.Println("Kasir API - Clean Architecture Implementation")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  api                Run the API server")
	fmt.Println("  api migrate        Run database migrations")
	fmt.Println("  api migrate-reset  Reset all migrations (drop all tables)")
	fmt.Println("  api seed           Seed database with sample data")
	fmt.Println("  api rls [on|off]   Enable or disable Row Level Security")
	fmt.Println("  api help           Show this help message")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help         Show this help message")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  APP_SERVER_HOST    Server host (default: localhost)")
	fmt.Println("  APP_SERVER_PORT    Server port (default: :8300)")
	fmt.Println("  APP_DATABASE_HOST  Database host")
	fmt.Println("  APP_DATABASE_USER  Database user")
	fmt.Println("  APP_DATABASE_PASSWORD  Database password")
	fmt.Println("  APP_DATABASE_DBNAME    Database name")
}

func runMigrations() {
	// Check for help flags
	if len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		fmt.Println("Run database migrations")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  api migrate        Run database migrations")
		fmt.Println()
		fmt.Println("Required Environment Variables:")
		fmt.Println("  APP_DATABASE_HOST      Database host")
		fmt.Println("  APP_DATABASE_USER      Database user")
		fmt.Println("  APP_DATABASE_PASSWORD  Database password")
		fmt.Println("  APP_DATABASE_DBNAME    Database name")
		return
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Database.Host == "" || cfg.Database.DBName == "" {
		logger.Error("Database configuration is required for migrations")
		log.Fatal("Database configuration is required for migrations")
	}

	db, err := database.NewPool(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Running database migrations...")
	if err := database.RunMigrations(db.DB, "database/migrations"); err != nil {
		logger.Error("Migration failed", "error", err)
		log.Fatalf("Migration failed: %v", err)
	}
	logger.Info("Migrations completed successfully")
}

func runMigrateReset() {
	// Check for help flags
	if len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		fmt.Println("Reset all database migrations")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  api migrate-reset  Reset migrations (runs all Down migrations)")
		fmt.Println()
		fmt.Println("Warning: This will drop all tables and data!")
		fmt.Println()
		fmt.Println("Required Environment Variables:")
		fmt.Println("  APP_DATABASE_HOST      Database host")
		fmt.Println("  APP_DATABASE_USER      Database user")
		fmt.Println("  APP_DATABASE_PASSWORD  Database password")
		fmt.Println("  APP_DATABASE_DBNAME    Database name")
		return
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Database.Host == "" || cfg.Database.DBName == "" {
		logger.Error("Database configuration is required for migrations")
		log.Fatal("Database configuration is required for migrations")
	}

	db, err := database.NewPool(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Warn("WARNING: This will drop all tables and data!")
	fmt.Print("Are you sure? Type 'yes' to confirm: ")

	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" {
		logger.Info("Migration reset cancelled")
		return
	}

	logger.Info("Resetting database migrations...")
	if err := database.ResetMigrations(db.DB, "database/migrations"); err != nil {
		logger.Error("Migration reset failed", "error", err)
		log.Fatalf("Migration reset failed: %v", err)
	}
	logger.Info("Migrations reset successfully")
}

func runSeeds() {
	// Check for help flags
	if len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		fmt.Println("Seed database with sample data")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  api seed           Seed database")
		fmt.Println()
		fmt.Println("Required Environment Variables:")
		fmt.Println("  APP_DATABASE_HOST      Database host")
		fmt.Println("  APP_DATABASE_USER      Database user")
		fmt.Println("  APP_DATABASE_PASSWORD  Database password")
		fmt.Println("  APP_DATABASE_DBNAME    Database name")
		return
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Database.Host == "" || cfg.Database.DBName == "" {
		logger.Error("Database configuration is required for seeding")
		log.Fatal("Database configuration is required for seeding")
	}

	db, err := database.NewPool(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Seeding database...")
	if err := database.RunSeeds(db.DB, "database/seeds"); err != nil {
		logger.Error("Seeding failed", "error", err)
		log.Fatalf("Seeding failed: %v", err)
	}
	logger.Info("Database seeded successfully")
}

func runRLS() {
	// Check for help flags
	if len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help") {
		fmt.Println("Enable or disable Row Level Security")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  api rls on         Enable RLS")
		fmt.Println("  api rls off        Disable RLS")
		fmt.Println()
		fmt.Println("Required Environment Variables:")
		fmt.Println("  APP_DATABASE_HOST      Database host")
		fmt.Println("  APP_DATABASE_USER      Database user")
		fmt.Println("  APP_DATABASE_PASSWORD  Database password")
		fmt.Println("  APP_DATABASE_DBNAME    Database name")
		return
	}

	if len(os.Args) < 3 {
		logger.Error("Usage: api rls [on|off]")
		log.Fatal("Usage: api rls [on|off]")
	}

	action := os.Args[2]
	if action != "on" && action != "off" {
		logger.Error("Usage: api rls [on|off]")
		log.Fatal("Usage: api rls [on|off]")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Database.Host == "" || cfg.Database.DBName == "" {
		logger.Error("Database configuration is required for RLS management")
		log.Fatal("Database configuration is required for RLS management")
	}

	db, err := database.NewPool(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if action == "on" {
		logger.Info("Enabling Row Level Security...")
		if err := database.EnableRLS(db.DB, "database/rls"); err != nil {
			logger.Error("Failed to enable RLS", "error", err)
			log.Fatalf("Failed to enable RLS: %v", err)
		}
		logger.Info("RLS enabled successfully")
	} else {
		logger.Info("Disabling Row Level Security...")
		if err := database.DisableRLS(db.DB, "database/rls"); err != nil {
			logger.Error("Failed to disable RLS", "error", err)
			log.Fatalf("Failed to disable RLS: %v", err)
		}
		logger.Info("RLS disabled successfully")
	}
}
