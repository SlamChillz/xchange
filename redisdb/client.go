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
	GetDel(context.Context, string) (string, error)
	Scan(context.Context, string, interface{}) error
	ScanDel(context.Context, string, interface{}) error
	Delete(context.Context, string) (int64, error)
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

func (rd *redisClient) GetDel(ctx context.Context, key string) (string, error) {
	storeValue, err := rd.client.GetDel(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return storeValue, nil
}

func (rd *redisClient) Scan(ctx context.Context, key string, val interface{}) error {
	return rd.client.Get(ctx, key).Scan(val)
}

func (rd *redisClient) ScanDel(ctx context.Context, key string, val interface{}) error {
	return rd.client.GetDel(ctx, key).Scan(val)
}

func (rd *redisClient) Delete(ctx context.Context, key string) (int64, error) {
	deleted, err := rd.client.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, err
}
