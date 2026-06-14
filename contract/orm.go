package contract

import (
	"context"

	"gorm.io/gorm"
)

// ORM is the port for GORM-based database access.
// Repositories depend on this interface, not on *gorm.DB directly.
type ORM interface {
	// Session returns a *gorm.DB scoped to the context.
	// Use for all queries: db.Session(ctx).Find(&list).Error
	Session(ctx context.Context) *gorm.DB

	// Transaction runs fn inside a database transaction, rolling back on error.
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	// Close releases the underlying connection pool.
	Close() error
}
