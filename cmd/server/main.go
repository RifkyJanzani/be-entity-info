package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rifky/be-entity-info/internal/config"
	"github.com/rifky/be-entity-info/internal/database"
	"github.com/rifky/be-entity-info/internal/handler"
	"github.com/rifky/be-entity-info/internal/repository"
	"github.com/rifky/be-entity-info/internal/router"
	"github.com/rifky/be-entity-info/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	// Run database migrations
	if err := database.RunMigrations(cfg.DatabaseURL(), "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Connect to PostgreSQL
	pool, err := database.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Wire dependencies
	repo := repository.NewEntityRepository(pool)
	svc := service.NewEntityService(repo)
	hdl := handler.NewEntityHandler(svc)

	// Setup router
	r := router.NewRouter(hdl, cfg.AllowedOrigins, cfg.AppEnv)

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}

	log.Println("Server stopped")
}
