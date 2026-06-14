package mock

import (
	"context"

	"github.com/bounkhongdev/kbgo/contract"
)

// Tx is a test double for contract.Tx.
// Override any Fn field to control behaviour per test.
// Default behaviour: Commit and Rollback succeed; queries return empty results.
type Tx struct {
	QueryFn    func(ctx context.Context, sql string, args ...any) (contract.Rows, error)
	QueryRowFn func(ctx context.Context, sql string, args ...any) contract.Row
	ExecFn     func(ctx context.Context, sql string, args ...any) error
	CommitFn   func(ctx context.Context) error
	RollbackFn func(ctx context.Context) error
}

func (t *Tx) Query(ctx context.Context, sql string, args ...any) (contract.Rows, error) {
	if t.QueryFn != nil {
		return t.QueryFn(ctx, sql, args...)
	}
	return &Rows{}, nil
}

func (t *Tx) QueryRow(ctx context.Context, sql string, args ...any) contract.Row {
	if t.QueryRowFn != nil {
		return t.QueryRowFn(ctx, sql, args...)
	}
	return &Row{}
}

func (t *Tx) Exec(ctx context.Context, sql string, args ...any) error {
	if t.ExecFn != nil {
		return t.ExecFn(ctx, sql, args...)
	}
	return nil
}

func (t *Tx) Commit(ctx context.Context) error {
	if t.CommitFn != nil {
		return t.CommitFn(ctx)
	}
	return nil
}

func (t *Tx) Rollback(ctx context.Context) error {
	if t.RollbackFn != nil {
		return t.RollbackFn(ctx)
	}
	return nil
}

// TransactionalDB wraps mock.Database and adds BeginTx support.
// Use this when your repository depends on contract.Transactional.
type TransactionalDB struct {
	Database
	BeginTxFn func(ctx context.Context) (contract.Tx, error)
}

func (d *TransactionalDB) BeginTx(ctx context.Context) (contract.Tx, error) {
	if d.BeginTxFn != nil {
		return d.BeginTxFn(ctx)
	}
	return &Tx{}, nil
}
