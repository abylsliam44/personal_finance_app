package redis_client

import (
	"github.com/go-redis/redis/v8"
)

// NewRedisClient создает и возвращает Redis клиент
func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
}