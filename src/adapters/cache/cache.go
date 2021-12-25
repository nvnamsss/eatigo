package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack"
)

const (
	shortTerm = 1000 * time.Millisecond
	midTerm   = 3000 * time.Millisecond
	longTerm  = 5000 * time.Millisecond
)

type CacheAdapter interface {
	Get(ctx context.Context, key string, v interface{}) error
	Set(ctx context.Context, key string, v interface{}, expiration time.Duration) error
}

type cachedAdapter struct {
	client *redis.Client
}

func NewCacheAdapter(client *redis.Client) CacheAdapter {
	return &cachedAdapter{
		client: client,
	}
}

func (r *cachedAdapter) Get(ctx context.Context, key string, v interface{}) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(shortTerm))
	defer cancel()
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return msgpack.Unmarshal(data, v)
}

func (r *cachedAdapter) Set(ctx context.Context, key string, v interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(midTerm))
	defer cancel()
	data, err := msgpack.Marshal(v)
	if err != nil {
		return nil
	}

	_, err = r.client.Set(ctx, key, data, expiration).Result()
	return err
}
