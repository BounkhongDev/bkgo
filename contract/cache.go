package contract

import (
	"context"
	"time"
)

// Cache is the port for any key-value cache adapter.
// Implement this interface to swap between Redis, Memcached, in-memory, etc.
//
// Note: Set JSON-encodes value before storing. Get returns the raw JSON string;
// callers must json.Unmarshal the result to recover the original value.
type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}
