package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Registry struct {
	conns map[string]*sql.DB
}

// OpenAll opens connections for provided DSNs keyed by name
func OpenAll(dbs map[string]string) (*Registry, error) {
	reg := &Registry{conns: make(map[string]*sql.DB)}
	for name, dsn := range dbs {
		conn, err := sql.Open("pgx", dsn)
		if err != nil {
			return nil, fmt.Errorf("open db %s: %w", name, err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := conn.PingContext(ctx); err != nil {
			cancel()
			_ = conn.Close()
			return nil, fmt.Errorf("ping db %s: %w", name, err)
		}
		cancel()
		reg.conns[name] = conn
	}
	return reg, nil
}

func (r *Registry) Get(name string) *sql.DB {
	return r.conns[name]
}

func (r *Registry) Close() error {
	var firstErr error
	for _, c := range r.conns {
		if err := c.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
