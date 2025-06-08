//go:generate swag init --generalInfo cmd/server/main.go --output ../../docs
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/Jexim/HelloGo/docs"
	mainMigration "github.com/Jexim/HelloGo/internal/repository/main/migrations"
	"github.com/Jexim/HelloGo/internal/rest"
	"github.com/Jexim/HelloGo/internal/rest/middleware"
	"github.com/Jexim/HelloGo/internal/svc/hello"
	"github.com/Jexim/HelloGo/pkg/config"
	"github.com/Jexim/HelloGo/pkg/logger"
	"github.com/Jexim/HelloGo/pkg/sentry"
)

// @title Hello Service API
// @version 1.0
// @description This is a hello service API documentation
// @host localhost:8080
// @BasePath /api
// @schemes http https
func main() {
	// Load config
	cfg := config.Load()

	// Init logger
	log := logger.New(viper.New())
	defer log.Sync()

	// Initialize Sentry
	if err := sentry.Init(cfg.Sentry.DSN, log); err != nil {
		log.Fatal("failed to initialize sentry", zap.Error(err))
	}
	defer sentry.Flush(2 * time.Second)

	// Create context that will be canceled on system interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Setup DB with retry logic
	db, err := setupDatabase(cfg.Database.URI, log)
	if err != nil {
		log.Fatal("failed to setup database", zap.Error(err))
	}

	// Setup HTTP server
	server, err := setupServer(cfg, db, log)
	if err != nil {
		log.Fatal("failed to setup server", zap.Error(err))
	}

	// Start server in a goroutine
	go func() {
		log.Info("starting server", zap.String("address", cfg.Server.Address))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info("shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	log.Info("server stopped")
}

func setupDatabase(dsn string, log *zap.Logger) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Warn("failed to connect to database, retrying...", zap.Error(err))
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after retries: %w", err)
	}

	// Automigrate
	if err := mainMigration.Migrate(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}

func setupServer(cfg *config.Config, db *gorm.DB, log *zap.Logger) (*http.Server, error) {
	mux := chi.NewRouter()

	// Middleware setup
	mux.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Trace-ID"},
		ExposedHeaders:   []string{"Link", "X-Trace-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	// Add trace ID middleware
	mux.Use(middleware.TraceID(log))

	// Metrics middleware
	mux.Use(middleware.Metrics)

	// Swagger documentation
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Metrics endpoint
	if cfg.Metrics.Enabled {
		mux.Handle(cfg.Metrics.Path, promhttp.Handler())
	}

	// Setup REST handlers
	_, err := rest.New(rest.InitArgs{
		Logger: log,
		DB:     db,
		Router: mux,
	}, rest.ArgsREST{
		Hello: hello.NewREST(mux, "/api/hello", hello.NewUsecase(hello.NewDatastore(db))),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create main REST: %w", err)
	}

	// Create server with timeouts
	return &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}, nil
}
