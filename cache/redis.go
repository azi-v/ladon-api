package cache

import (
	"git.ymt360.com/usercenter/ymt-ladon/config"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func Conn(cfg *config.CacheConfig) {
	rdb =redis.NewClient(&redis.Options{
		Addr:"",
		Password: "",
		DB:0,
	})
}
