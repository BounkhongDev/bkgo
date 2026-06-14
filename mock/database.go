package mock

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bounkhongdev/kbgo/contract"
)

// Database is a test double for contract.Database.
// Override any Fn field to control what the mock returns.
type Database struct {
	QueryFn    func(ctx context.Context, sql string, args ...any) (contract.Rows, error)
	QueryRowFn func(ctx context.Context, sql string, args ...any) contract.Row
	ExecFn     func(ctx context.Context, sql string, args ...any) error
	PingFn     func(ctx context.Context) error
	CloseFn    func()
}

func (m *Database) Query(ctx context.Context, sql string, args ...any) (contract.Rows, error) {
	if m.QueryFn != nil {
		return m.QueryFn(ctx, sql, args...)
	}
	return &Rows{}, nil
}

func (m *Database) QueryRow(ctx context.Context, sql string, args ...any) contract.Row {
	if m.QueryRowFn != nil {
		return m.QueryRowFn(ctx, sql, args...)
	}
	return &Row{}
}

func (m *Database) Exec(ctx context.Context, sql string, args ...any) error {
	if m.ExecFn != nil {
		return m.ExecFn(ctx, sql, args...)
	}
	return nil
}

func (m *Database) Ping(ctx context.Context) error {
	if m.PingFn != nil {
		return m.PingFn(ctx)
	}
	return nil
}

func (m *Database) Close() {
	if m.CloseFn != nil {
		m.CloseFn()
	}
}

// ── Row helpers ───────────────────────────────────────────────────────────────

// Row implements contract.Row for testing.
// Use NewRow(...values) to create one that Scan() reads from.
type Row struct {
	values []any
	err    error
}

// NewRow creates a mock Row whose Scan fills destinations in order.
func NewRow(values ...any) *Row { return &Row{values: values} }

// NewRowError creates a mock Row that always returns err from Scan.
func NewRowError(err error) *Row { return &Row{err: err} }

func (r *Row) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) > len(r.values) {
		return fmt.Errorf("mock: Scan wants %d values, row has %d", len(dest), len(r.values))
	}
	for i, d := range dest {
		reflect.ValueOf(d).Elem().Set(reflect.ValueOf(r.values[i]))
	}
	return nil
}

// ── Rows helpers ──────────────────────────────────────────────────────────────

// Rows implements contract.Rows for testing.
// Use NewRows([][]any{...}) to create one backed by a slice of value sets.
type Rows struct {
	data [][]any
	idx  int
	err  error
}

// NewRows creates mock Rows from a 2-D slice of values.
func NewRows(data [][]any) *Rows { return &Rows{data: data, idx: -1} }

// NewRowsError creates mock Rows whose Err() returns the given error.
func NewRowsError(err error) *Rows { return &Rows{err: err} }

func (r *Rows) Next() bool {
	r.idx++
	return r.idx < len(r.data)
}

func (r *Rows) Scan(dest ...any) error {
	row := r.data[r.idx]
	if len(dest) > len(row) {
		return fmt.Errorf("mock: Scan wants %d values, row has %d", len(dest), len(row))
	}
	for i, d := range dest {
		reflect.ValueOf(d).Elem().Set(reflect.ValueOf(row[i]))
	}
	return nil
}

func (r *Rows) Close() {}
func (r *Rows) Err() error { return r.err }
