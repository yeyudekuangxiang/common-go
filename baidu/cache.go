package baidu

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"time"
)

type ICache interface {
	GetValue(key string) (value string, err error)
	SetValue(key string, value string, expiration time.Duration) error
}

func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{
		Client: client,
		Prefix: prefix,
	}
}

type RedisCache struct {
	Client *redis.Client
	Prefix string
}

func (c *RedisCache) GetValue(key string) (value string, err error) {
	val, err := c.Client.Get(context.Background(), c.formatKey(key)).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return val, nil
}
func (c *RedisCache) SetValue(key string, value string, expiration time.Duration) error {
	return c.Client.Set(context.Background(), c.formatKey(key), value, expiration).Err()
}
func (c *RedisCache) formatKey(key string) string {
	return c.Prefix + key
}

func NewMemoryCache(prefix string) *MemoryCache {
	return &MemoryCache{
		cache:  cache.New(time.Minute*5, time.Minute*5),
		Prefix: prefix,
	}
}

type MemoryCache struct {
	cache  *cache.Cache
	Prefix string
}

func (c *MemoryCache) GetValue(key string) (value string, err error) {
	val, ok := c.cache.Get(c.formatKey(key))
	if ok {
		return val.(string), nil
	}
	return "", nil
}
func (c *MemoryCache) SetValue(key string, value string, expiration time.Duration) error {
	c.cache.Set(c.formatKey(key), value, expiration)
	return nil
}
func (c *MemoryCache) formatKey(key string) string {
	return c.Prefix + key
}
