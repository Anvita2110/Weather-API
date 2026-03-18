package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var globalRedisClient *RedisClient

func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	globalRedisClient = &RedisClient{client: rdb}
	return globalRedisClient
}

func CloseRedisClient() {
	if globalRedisClient != nil {
		globalRedisClient.client.Close()
	}
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, json, expiration).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
