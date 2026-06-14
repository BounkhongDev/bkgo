package postgres

import (
	"context"
	"fmt"

	"github.com/BounkhongDev/bkgo/contract"
	"github.com/jackc/pgx/v5"
)

// BeginTx starts a new transaction. Satisfies contract.Transactional.
func (d *DB) BeginTx(ctx context.Context) (contract.Tx, error) {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres: begin tx: %w", err)
	}
	return &pgTx{tx: tx}, nil
}

type pgTx struct{ tx pgx.Tx }

func (t *pgTx) Query(ctx context.Context, sql string, args ...any) (contract.Rows, error) {
	rows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgRows{rows: rows}, nil
}

func (t *pgTx) QueryRow(ctx context.Context, sql string, args ...any) contract.Row {
	return &pgRow{row: t.tx.QueryRow(ctx, sql, args...)}
}

func (t *pgTx) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := t.tx.Exec(ctx, sql, args...)
	return err
}

func (t *pgTx) Commit(ctx context.Context) error   { return t.tx.Commit(ctx) }
func (t *pgTx) Rollback(ctx context.Context) error { return t.tx.Rollback(ctx) }
