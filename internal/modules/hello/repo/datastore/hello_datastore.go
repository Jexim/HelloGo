package datastore

import (
	"context"
	"database/sql"

	"github.com/Jexim/HelloGo/internal/modules/hello/model"
	gen "github.com/Jexim/HelloGo/internal/modules/hello/repo/sqlc/gen"
)

type helloDatastore struct {
	db *sql.DB
	q  *gen.Queries
}

func NewDatastore(db *sql.DB) model.Datastore {
	return &helloDatastore{db: db, q: gen.New(db)}
}

// Create creates a new hello in the database
func (d *helloDatastore) Create(ctx context.Context, in *model.Hello) (*model.Hello, error) {
	h, err := d.q.CreateHello(ctx, in.Message)
	if err != nil {
		return nil, err
	}
	return &model.Hello{ID: uint(h.ID), Message: h.Message}, nil
}

// GetAll retrieves all hellos from the database
func (d *helloDatastore) GetAll(ctx context.Context, limit, offset int) ([]model.Hello, error) {
	list, err := d.q.ListHellos(ctx, gen.ListHellosParams{Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return nil, err
	}
	result := make([]model.Hello, 0, len(list))
	for _, it := range list {
		result = append(result, model.Hello{ID: uint(it.ID), Message: it.Message})
	}
	return result, nil
}

// Get retrieves a hello by ID from the database
func (d *helloDatastore) Get(ctx context.Context, id int) (*model.Hello, error) {
	h, err := d.q.GetHello(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return &model.Hello{ID: uint(h.ID), Message: h.Message}, nil
}

// Update updates a hello in the database
func (d *helloDatastore) Update(ctx context.Context, id int, in *model.Hello) error {
	return d.q.UpdateHello(ctx, gen.UpdateHelloParams{Message: in.Message, ID: int32(id)})
}

// Delete removes a hello from the database
func (d *helloDatastore) Delete(ctx context.Context, id int) error {
	return d.q.DeleteHello(ctx, int32(id))
}
