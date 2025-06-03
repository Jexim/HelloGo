package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type QueryArgs map[string]any

type Base struct {
	// DB is the GORM database instance
	DB *gorm.DB
}

func NewBase(db *gorm.DB) Base {
	return Base{
		DB: db,
	}
}

type ExtendedFilter struct {
	Filter    *Filter
	FromEvent *EventIdentity
	ToEvent   *EventIdentity
	Search    *string
	TimeFrom  *time.Time
	TimeTo    *time.Time
}

func (f *ExtendedFilter) SearchNotEmpty() bool {
	return f.Search != nil && len(*f.Search) != 0
}

// Filter represents a basic filter structure
type Filter struct {
	Limit  int
	Offset int
}

// EventIdentity represents an event identifier
type EventIdentity struct {
	ID   uint
	Name string
}

// WithContext returns a new GORM DB instance with the given context
func (b *Base) WithContext(ctx context.Context) *gorm.DB {
	return b.DB.WithContext(ctx)
}

// Transaction executes the given function within a transaction
func (b *Base) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return b.DB.WithContext(ctx).Transaction(fn)
}
