package contract

import "context"

// Tx represents an active database transaction.
type Tx interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, args ...any) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Transactional extends Database with transaction support.
// The postgres adapter implements this; use it when your usecase needs atomicity.
type Transactional interface {
	Database
	BeginTx(ctx context.Context) (Tx, error)
}
