package hello

import (
	"database/sql"

	"github.com/go-chi/chi"

	"github.com/Jexim/HelloGo/internal/modules/hello/model"
	"github.com/Jexim/HelloGo/internal/modules/hello/repo/datastore"
	"github.com/Jexim/HelloGo/internal/modules/hello/rest"
	"github.com/Jexim/HelloGo/internal/modules/hello/usecase"
)

type (
	Datastore = model.Datastore

	Hello = model.Hello

	Usecase = model.Usecase

	RESTHello = model.REST
)

func NewDatastore(db *sql.DB) Datastore {
	return datastore.NewDatastore(db)
}

func NewUsecase(ds Datastore) Usecase {
	return usecase.New(ds)
}

func NewREST(mux *chi.Mux, prefix string, helloUC Usecase) RESTHello {
	return rest.New(mux, prefix, helloUC)
}

var (
	ErrNotFound = model.ErrNotFound
)
