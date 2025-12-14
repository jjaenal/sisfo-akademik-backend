package redisutil

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type fakeKV struct {
	store map[string]string
}

func (f *fakeKV) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	if f.store == nil {
		f.store = map[string]string{}
	}
	f.store[key] = value.(string)
	cmd := redis.NewStatusCmd(ctx)
	return cmd
}
func (f *fakeKV) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	cmd.SetVal(f.store[key])
	return cmd
}
func (f *fakeKV) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(1)
	return cmd
}

type fakeCounter struct {
	count map[string]int64
	expireCalls int
}

func (f *fakeCounter) Incr(ctx context.Context, key string) *redis.IntCmd {
	if f.count == nil {
		f.count = map[string]int64{}
	}
	f.count[key]++
	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(f.count[key])
	return cmd
}
func (f *fakeCounter) Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd {
	f.expireCalls++
	cmd := redis.NewBoolCmd(ctx)
	cmd.SetVal(true)
	return cmd
}

func TestKV(t *testing.T) {
	var kv KV = &fakeKV{}
	ctx := context.Background()
	if err := Set(ctx, kv, "k", "v", time.Minute); err != nil {
		t.Fatalf("Set error: %v", err)
	}
	val, err := Get(ctx, kv, "k")
	if err != nil || val != "v" {
		t.Fatalf("Get mismatch: %v %s", err, val)
	}
	if err := Del(ctx, kv, "k"); err != nil {
		t.Fatalf("Del error: %v", err)
	}
}

type fakeKVE struct{}

func (f *fakeKVE) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	cmd.SetErr(errors.New("set failed"))
	return cmd
}
func (f *fakeKVE) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	cmd.SetErr(errors.New("get failed"))
	return cmd
}
func (f *fakeKVE) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	cmd.SetErr(errors.New("del failed"))
	return cmd
}

func TestKVErrors(t *testing.T) {
	ctx := context.Background()
	var kv KV = &fakeKVE{}
	if err := Set(ctx, kv, "k", "v", time.Minute); err == nil {
		t.Fatalf("expected error from Set")
	}
	if _, err := Get(ctx, kv, "k"); err == nil {
		t.Fatalf("expected error from Get")
	}
	if err := Del(ctx, kv, "k"); err == nil {
		t.Fatalf("expected error from Del")
	}
}

func TestLimiterAdapter(t *testing.T) {
	c := &fakeCounter{}
	lim := NewLimiterFromCounter(c)
	ctx := context.Background()
	n, err := lim.Incr(ctx, "x")
	if err != nil || n != 1 {
		t.Fatalf("Incr failed")
	}
	if err := lim.Expire(ctx, "x", time.Minute); err != nil {
		t.Fatalf("Expire failed")
	}
	// Pastikan expire dipanggil sesuai skenario pertama
	if c.expireCalls == 0 {
		t.Fatalf("expected expire to be called at least once")
	}
	// Panggilan kedua tidak seharusnya memicu expire lagi pada helper IncrWithTTL
	n2, err := lim.Incr(ctx, "x")
	if err != nil || n2 != 2 {
		t.Fatalf("second Incr failed")
	}
}

func TestIncrWithTTLSuccess(t *testing.T) {
	c := &fakeCounter{}
	ctx := context.Background()
	n1, err := IncrWithTTL(ctx, c, "k", time.Minute)
	if err != nil || n1 != 1 {
		t.Fatalf("first IncrWithTTL failed")
	}
	if c.expireCalls != 1 {
		t.Fatalf("expire should be called on first increment")
	}
	n2, err := IncrWithTTL(ctx, c, "k", time.Minute)
	if err != nil || n2 != 2 {
		t.Fatalf("second IncrWithTTL failed")
	}
	if c.expireCalls != 1 {
		t.Fatalf("expire should not be called again")
	}
}

func TestClientNewRaw(t *testing.T) {
	cl := New("localhost:6379")
	if cl == nil || cl.Raw() == nil {
		t.Fatalf("client or raw should not be nil")
	}
}

type errorCounter struct{}

func (e *errorCounter) Incr(ctx context.Context, key string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	cmd.SetErr(errors.New("incr failed"))
	return cmd
}
func (e *errorCounter) Expire(ctx context.Context, key string, ttl time.Duration) *redis.BoolCmd {
	cmd := redis.NewBoolCmd(ctx)
	cmd.SetVal(false)
	return cmd
}

func TestIncrWithTTLError(t *testing.T) {
	ctx := context.Background()
	_, err := IncrWithTTL(ctx, &errorCounter{}, "k", time.Minute)
	if err == nil {
		t.Fatalf("expected error from IncrWithTTL")
	}
}
