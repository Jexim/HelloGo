package rest

import (
	"go.uber.org/zap"

	"github.com/Jexim/HelloGo/internal/svc/hello"
)

type REST struct {
	Hello hello.RESTHello

	logger *zap.Logger
}

type InitArgs struct {
	Logger *zap.Logger
}

type ArgsREST struct {
	Hello hello.RESTHello
}

func New(args InitArgs, argsREST ArgsREST) (*REST, error) {
	return &REST{
		logger: args.Logger,
		Hello:  argsREST.Hello,
	}, nil
}
