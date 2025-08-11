//go:generate swag init --generalInfo cmd/server/main.go --output ../../docs
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/Jexim/HelloGo/docs"
	httpadapter "github.com/Jexim/HelloGo/internal/adapter/http"
	httpmw "github.com/Jexim/HelloGo/internal/adapter/http/middleware"

	// restrespond "github.com/Jexim/HelloGo/internal/rest/respond"
	"github.com/Jexim/HelloGo/internal/modules/hello"
	"github.com/Jexim/HelloGo/internal/platform/config"
	platformdb "github.com/Jexim/HelloGo/internal/platform/db"
	"github.com/Jexim/HelloGo/internal/platform/logger"
	"github.com/Jexim/HelloGo/internal/platform/sentry"
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
	log := logger.New(cfg.Logger.Level)
	defer log.Sync()

	// Initialize Sentry
	if err := sentry.Init(cfg.Sentry.DSN, cfg.Sentry.Environment, log); err != nil {
		log.Fatal("failed to initialize sentry", zap.Error(err))
	}
	defer sentry.Flush(2 * time.Second)

	// Create context that will be canceled on system interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Setup multiple DBs via registry
	mainDB, reg, err := setupDatabases(cfg, log)
	if err != nil {
		log.Fatal("failed to setup database(s)", zap.Error(err))
	}
	defer reg.Close()

	// Setup HTTP server
	server, err := setupServer(cfg, mainDB, log)
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

func setupDatabases(cfg *config.Config, log *zap.Logger) (*sql.DB, *platformdb.Registry, error) {
	dsnMap := make(map[string]string, len(cfg.Databases))
	for name, dc := range cfg.Databases {
		dsnMap[name] = dc.URI
	}
	reg, err := platformdb.OpenAll(dsnMap)
	if err != nil {
		return nil, nil, err
	}
	mainDB := reg.Get("main")
	if mainDB == nil {
		return nil, reg, fmt.Errorf("main database is not configured")
	}
	return mainDB, reg, nil
}

func setupServer(cfg *config.Config, db *sql.DB, log *zap.Logger) (*http.Server, error) {
	mux := chi.NewRouter()

	// Middleware setup
	mux.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Trace-ID"},
		ExposedHeaders:   []string{"Link", "X-Trace-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}).Handler)

	// Add trace ID middleware
	mux.Use(httpmw.TraceID(log))

	// Metrics middleware
	mux.Use(httpmw.Metrics)

	// Error middleware with Sentry capture
	mux.Use(httpmw.ErrorHandler(log, func(err error) {
		// Capture to Sentry (if DSN configured)
		if cfg.Sentry.DSN != "" {
			sentry.Recover(log) // ensures panic capture; for non-panics, send as exception
		}
	}))

	// Swagger documentation
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Metrics endpoint
	if cfg.Metrics.Enabled {
		mux.Handle(cfg.Metrics.Path, promhttp.Handler())
	}

	// Setup REST handlers
	_, err := httpadapter.New(httpadapter.InitArgs{
		Logger: log,
		DB:     db,
		Router: mux,
	}, httpadapter.ArgsREST{
		Hello: hello.NewREST(mux, "/api/v1/hello", hello.NewUsecase(hello.NewDatastore(db))),
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
