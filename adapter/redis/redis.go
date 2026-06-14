package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bounkhongdev/kbgo/config"
	goredis "github.com/redis/go-redis/v9"
)

// Cache is the Redis adapter implementing contract.Cache.
type Cache struct {
	client *goredis.Client
}

// New creates a Redis client and verifies connectivity.
func New(ctx context.Context, cfg config.Redis) (*Cache, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: connect: %w", err)
	}
	return &Cache{client: client}, nil
}

// Client exposes the underlying goredis.Client for advanced usage.
func (c *Cache) Client() *goredis.Client { return c.client }

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis: marshal: %w", err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (c *Cache) Close() error { return c.client.Close() }
