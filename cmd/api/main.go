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
)

func main() {
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

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize repositories
	var productRepo repository.ProductReader
	var productWriter repository.ProductWriter
	var categoryRepo repository.CategoryReader
	var categoryWriter repository.CategoryWriter
	var db *database.DB

	// Check if database is configured
	if cfg.Database.Host != "" && cfg.Database.DBName != "" {
		// Use PostgreSQL
		db, err = database.NewPool(cfg.Database)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		// Test connection
		if err := db.Ping(context.Background()); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		fmt.Println("Connected to PostgreSQL database")

		// Create PostgreSQL repositories
		pgProductRepo := postgres.NewProductRepository(db.DB)
		productRepo = pgProductRepo
		productWriter = pgProductRepo

		pgCategoryRepo := postgres.NewCategoryRepository(db.DB)
		categoryRepo = pgCategoryRepo
		categoryWriter = pgCategoryRepo
	} else {
		// Use in-memory repositories
		fmt.Println("Using in-memory storage")
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

	// Initialize handlers
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	var healthHandler *handler.HealthHandler
	if db != nil {
		healthHandler = handler.NewHealthHandler(db.DB)
	} else {
		healthHandler = handler.NewHealthHandler(nil)
	}

	// Setup routes
	mux := http.NewServeMux()
	handler.SetupRoutes(mux, productHandler, categoryHandler, healthHandler)

	// Create server
	server := &http.Server{
		Addr:         cfg.Server.Host + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("Starting server on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Fatalf("Server failed: %v", err)
	case sig := <-quit:
		fmt.Printf("Received signal: %v\n", sig)
		fmt.Println("Shutting down server...")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		fmt.Println("Server stopped gracefully")
	}
}

func printHelp() {
	fmt.Println("Kasir API - Clean Architecture Implementation")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  api                Run the API server")
	fmt.Println("  api migrate        Run database migrations")
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
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Database.Host == "" || cfg.Database.DBName == "" {
		log.Fatal("Database configuration is required for migrations")
	}

	db, err := database.NewPool(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Running database migrations...")
	if err := database.RunMigrations(db.DB, "migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations completed successfully")
}
