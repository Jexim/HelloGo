//go:generate swag init --generalInfo cmd/server/main.go --output ../../docs
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/Jexim/HelloGo/docs"
	mainMigration "github.com/Jexim/HelloGo/internal/repository/main/migrations"
	"github.com/Jexim/HelloGo/internal/rest"
	"github.com/Jexim/HelloGo/internal/svc/hello"
	"github.com/Jexim/HelloGo/pkg/config"
	"github.com/Jexim/HelloGo/pkg/logger"
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

	// Create context that will be canceled on system interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Setup DB with retry logic
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(cfg.Database.URI), &gorm.Config{})
		if err == nil {
			break
		}
		log.Warn("failed to connect to database, retrying...", zap.Error(err))
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Fatal("failed to connect to database after retries", zap.Error(err))
	}

	// Automigrate
	if err := mainMigration.Migrate(db); err != nil {
		log.Fatal("Migration failed", zap.Error(err))
	}

	mux := chi.NewRouter()

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	mux.Use(corsMiddleware.Handler)

	// Swagger documentation
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	_, err = rest.New(rest.InitArgs{
		Logger: log,
	}, rest.ArgsREST{
		Hello: hello.NewREST(mux, "/api/hello", hello.NewUsecase(hello.NewDatastore(db))),
	})
	if err != nil {
		log.Fatal("failed to create main REST", zap.Error(err))
	}

	// Create server with timeouts
	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
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
