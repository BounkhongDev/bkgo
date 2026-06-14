package contract

import "context"

// Row represents a single query result row.
type Row interface {
	Scan(dest ...any) error
}

// Rows represents multiple query result rows.
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

// Database is the port for any relational database adapter.
// Implement this interface to swap between PostgreSQL, MySQL, SQLite, etc.
type Database interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, args ...any) error
	Ping(ctx context.Context) error
	Close()
}
