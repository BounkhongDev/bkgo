package mock

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// Cache is a thread-safe in-memory implementation of contract.Cache for testing.
// Values are JSON-marshalled on Set — identical to the real Redis adapter —
// so tests accurately reflect production behaviour.
// TTL is accepted but ignored (no expiry in tests).
type Cache struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewCache() *Cache {
	return &Cache{store: make(map[string]string)}
}

func (c *Cache) Set(_ context.Context, key string, value any, _ time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = string(data)
	return nil
}

func (c *Cache) Get(_ context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.store[key]
	if !ok {
		return "", errors.New("mock: key not found")
	}
	return v, nil
}

func (c *Cache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	return nil
}

func (c *Cache) Exists(_ context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.store[key]
	return ok, nil
}

func (c *Cache) Close() error { return nil }

// Flush clears all keys — call between test cases to reset state.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]string)
}
