package mock

import (
	"context"

	"gorm.io/gorm"
)

// ORM is a test double for contract.ORM.
// Set SessionFn and TransactionFn to control behaviour per test.
// Unset fields panic if called — tests fail loudly on unexpected calls.
type ORM struct {
	SessionFn     func(ctx context.Context) *gorm.DB
	TransactionFn func(ctx context.Context, fn func(tx *gorm.DB) error) error
	CloseFn       func() error
}

func (m *ORM) Session(ctx context.Context) *gorm.DB {
	if m.SessionFn == nil {
		panic("mock.ORM: SessionFn is not set")
	}
	return m.SessionFn(ctx)
}

func (m *ORM) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	if m.TransactionFn != nil {
		return m.TransactionFn(ctx, fn)
	}
	if m.SessionFn == nil {
		panic("mock.ORM: SessionFn is not set (required by default Transaction fallback)")
	}
	return fn(m.SessionFn(ctx))
}

func (m *ORM) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	return nil
}
