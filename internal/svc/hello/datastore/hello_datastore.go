package datastore

import (
	"context"

	"gorm.io/gorm"

	"github.com/Jexim/HelloGo/internal/repository"
	"github.com/Jexim/HelloGo/internal/svc/hello/model"
)

type helloDatastore struct {
	repository.Base
}

func NewDatastore(db *gorm.DB) model.Datastore {
	return &helloDatastore{
		Base: repository.NewBase(db),
	}
}

type hello struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Message string `json:"message"`
}

// Create creates a new hello in the database
func (d *helloDatastore) Create(ctx context.Context, in *model.Hello) (*model.Hello, error) {
	hello := &hello{
		Message: in.Message,
	}

	err := d.WithContext(ctx).Create(hello).Error
	if err != nil {
		return nil, err
	}

	return &model.Hello{
		ID:      hello.ID,
		Message: hello.Message,
	}, nil
}

// GetAll retrieves all hellos from the database
func (d *helloDatastore) GetAll(ctx context.Context) ([]model.Hello, error) {
	var rows []hello
	err := d.WithContext(ctx).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	hellos := make([]model.Hello, len(rows))
	for i, row := range rows {
		hellos[i] = model.Hello{
			Message: row.Message,
		}
	}
	return hellos, nil
}

// Get retrieves a hello by ID from the database
func (d *helloDatastore) Get(ctx context.Context, id int) (*model.Hello, error) {
	var row hello
	err := d.WithContext(ctx).First(&row, id).Error
	if err != nil {
		return nil, err
	}

	return &model.Hello{
		ID:      row.ID,
		Message: row.Message,
	}, nil
}

// Update updates a hello in the database
func (d *helloDatastore) Update(ctx context.Context, id int, in *model.Hello) error {
	row := hello{
		ID:      uint(id),
		Message: in.Message,
	}
	return d.WithContext(ctx).Save(&row).Error
}

// Delete removes a hello from the database
func (d *helloDatastore) Delete(ctx context.Context, id int) error {
	return d.WithContext(ctx).Delete(&hello{}, id).Error
}
