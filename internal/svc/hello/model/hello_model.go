package model

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrNotFound = errors.New("hello not found")
)

//go:generate go run -mod=mod go.uber.org/mock/mockgen -mock_names Datastore=MockedDatastore -package mock -destination ../mock/hello_datastore_mock.go . Datastore
type Datastore interface {
	Create(ctx context.Context, hello *Hello) (*Hello, error)
	GetAll(ctx context.Context) ([]Hello, error)
	Get(ctx context.Context, id int) (*Hello, error)
	Update(ctx context.Context, id int, hello *Hello) error
	Delete(ctx context.Context, id int) error
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -mock_names Usecase=MockedUsecase -package mock -destination ../mock/hello_usecase_mock.go . Usecase
type Usecase interface {
	Datastore
}

type REST interface {
	GetHello(w http.ResponseWriter, r *http.Request)
}

// Hello represents
// @Description Hello
type Hello struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
