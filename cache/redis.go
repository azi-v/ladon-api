package cache

import (
	"github.com/azi-v/ladon-api/config"
	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func Conn(cfg *config.RedisConfig) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       int(cfg.DB),
		PoolSize: int(cfg.Poolsize),
	})
}
