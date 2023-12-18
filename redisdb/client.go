package redisdb

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrKeyNotFound = redis.Nil
	ErrorKeyNotSet = errors.New("key not set")
)

type RedisClient interface {
	Set(context.Context, string, interface{}, time.Duration) (string, error)
	Get(context.Context, string) (string, error)
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) RedisClient {
	return &redisClient{
		client: client,
	}
}

func (rd *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	storeValue, err := rd.client.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return "", err
	}
	return storeValue, nil
}

func (rd *redisClient) Get(ctx context.Context, key string) (string, error) {
	storeValue, err := rd.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return storeValue, nil
}
