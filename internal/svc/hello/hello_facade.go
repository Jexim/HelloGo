package hello

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"

	"github.com/Jexim/HelloGo/internal/svc/hello/datastore"
	"github.com/Jexim/HelloGo/internal/svc/hello/model"
	"github.com/Jexim/HelloGo/internal/svc/hello/rest"
	"github.com/Jexim/HelloGo/internal/svc/hello/usecase"
)

type (
	Datastore = model.Datastore

	Hello = model.Hello

	Usecase = model.Usecase

	RESTHello = model.REST
)

func NewDatastore(db *gorm.DB) Datastore {
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
