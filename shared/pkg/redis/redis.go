package redisutil

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter interface {
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

type Counter interface {
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd
}

type KV interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type Client struct {
	raw *redis.Client
}

func New(addr string) *Client {
	return &Client{
		raw: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (c *Client) Raw() *redis.Client {
	return c.raw
}

type limiterAdapter struct {
	c Counter
}

func NewLimiterFromCounter(c Counter) Limiter {
	return &limiterAdapter{c: c}
}

func (a *limiterAdapter) Incr(ctx context.Context, key string) (int64, error) {
	return a.c.Incr(ctx, key).Result()
}
func (a *limiterAdapter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return a.c.Expire(ctx, key, ttl).Err()
}

func Set(ctx context.Context, k KV, key string, value any, ttl time.Duration) error {
	return k.Set(ctx, key, value, ttl).Err()
}

func Get(ctx context.Context, k KV, key string) (string, error) {
	return k.Get(ctx, key).Result()
}

func Del(ctx context.Context, k KV, key string) error {
	return k.Del(ctx, key).Err()
}

func IncrWithTTL(ctx context.Context, c Counter, key string, ttl time.Duration) (int64, error) {
	n, err := c.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if n == 1 {
		_ = c.Expire(ctx, key, ttl).Err()
	}
	return n, nil
}
