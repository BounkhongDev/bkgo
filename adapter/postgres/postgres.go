package postgres

import (
	"context"
	"fmt"

	"github.com/bounkhongdev/kbgo/config"
	"github.com/bounkhongdev/kbgo/contract"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the PostgreSQL adapter implementing contract.Database.
type DB struct {
	pool *pgxpool.Pool
}

// New creates a new PostgreSQL connection pool and verifies connectivity.
func New(ctx context.Context, cfg config.Postgres) (*DB, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("postgres: connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}
	return &DB{pool: pool}, nil
}

// Pool exposes the underlying pgxpool.Pool for advanced usage (e.g. transactions).
func (d *DB) Pool() *pgxpool.Pool { return d.pool }

func (d *DB) Query(ctx context.Context, sql string, args ...any) (contract.Rows, error) {
	rows, err := d.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgRows{rows: rows}, nil
}

func (d *DB) QueryRow(ctx context.Context, sql string, args ...any) contract.Row {
	return &pgRow{row: d.pool.QueryRow(ctx, sql, args...)}
}

func (d *DB) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := d.pool.Exec(ctx, sql, args...)
	return err
}

func (d *DB) Ping(ctx context.Context) error { return d.pool.Ping(ctx) }
func (d *DB) Close()                         { d.pool.Close() }

// pgRows wraps pgx.Rows to satisfy contract.Rows.
type pgRows struct{ rows pgx.Rows }

func (r *pgRows) Next() bool            { return r.rows.Next() }
func (r *pgRows) Scan(dest ...any) error { return r.rows.Scan(dest...) }
func (r *pgRows) Close()                { r.rows.Close() }
func (r *pgRows) Err() error            { return r.rows.Err() }

// pgRow wraps pgx.Row to satisfy contract.Row.
type pgRow struct{ row pgx.Row }

func (r *pgRow) Scan(dest ...any) error { return r.row.Scan(dest...) }
