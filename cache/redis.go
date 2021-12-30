package cache

import (
	"git.ymt360.com/usercenter/ymt-ladon/config"
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
