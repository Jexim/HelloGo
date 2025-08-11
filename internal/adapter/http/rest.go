package httpadapter

import (
	"database/sql"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	healthrest "github.com/Jexim/HelloGo/internal/adapter/http/health"
	"github.com/Jexim/HelloGo/internal/modules/hello"
	healthcheck "github.com/Jexim/HelloGo/internal/platform/health"
)

type REST struct {
	Hello  hello.RESTHello
	Health *healthrest.REST

	logger *zap.Logger
	router chi.Router
}

type InitArgs struct {
	Logger *zap.Logger
	DB     *sql.DB
	Router chi.Router
}

type ArgsREST struct {
	Hello hello.RESTHello
}

func New(args InitArgs, argsREST ArgsREST) (*REST, error) {
	// Initialize health checker
	healthChecker := healthcheck.NewChecker(args.DB, args.Logger)

	return &REST{
		logger: args.Logger,
		Hello:  argsREST.Hello,
		Health: healthrest.New(args.Router, "/health", healthChecker, args.Logger),
		router: args.Router,
	}, nil
}
