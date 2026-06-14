package gormadapter

import (
	"context"
	"fmt"

	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/BounkhongDev/bkgo/config"
)

// DB wraps *gorm.DB and satisfies contract.ORM.
type DB struct {
	db *gorm.DB
}

// New opens a GORM PostgreSQL connection pool using the provided config.
func New(cfg config.Postgres) (*DB, error) {
	db, err := gorm.Open(gormpg.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gorm.DB(): %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	return &DB{db: db}, nil
}

// Session returns a *gorm.DB scoped to the given context.
func (d *DB) Session(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

// Transaction runs fn inside a database transaction, rolling back on error.
func (d *DB) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn)
}

// Close releases the underlying connection pool.
func (d *DB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Raw exposes the underlying *gorm.DB for advanced use (AutoMigrate, etc.).
func (d *DB) Raw() *gorm.DB {
	return d.db
}
