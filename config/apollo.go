package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log"
)

func InitApolloConfig(ctx context.Context) *ApolloConfig {
	// 获取env apollo的token等
	res, err := http.Get("http://dev-apollo-proxy.ymt.io/configs/gogodog/default/application?&token=LRvFPj91aYPgsMcfJJxZl6blbKm9ADGW")
	if err != nil {
		log.Fatal(ctx, err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatal(ctx, "Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(ctx, err)
	}
	fmt.Printf("%s", body)

	apolloResp := &ApolloConfig{}
	err = json.Unmarshal(body, apolloResp)
	if err != nil {
		log.Fatal(ctx, err)
	}

	return apolloResp
}

type ApolloConfig struct {
	LogConfig   *LogConfig
	CacheConfig *RedisConfig
	DBConfig    *MongoConfig
}

type LogConfig struct {
}

type RedisConfig struct {
	Addr     string `json:"addr,omitempty"`
	Password string `json:"password,omitempty"`
	DB       int32  `json:"db,omitempty"`
	Poolsize int32  `json:"poolsize,omitempty"`
}



type MongoConfig struct {
	User        string `json:"user,omitempty"`
	Password    string `json:"password,omitempty"`
	Host        string `json:"host,omitempty"`
	Port        int32  `json:"port,omitempty"`
	MaxPoolSize int32  `json:"max_pool_size,omitempty"`
	W           string `json:"w,omitempty"`
}

func (mc *MongoConfig) GetURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/?maxPoolSize=%d&w=%s", mc.User, mc.Password, mc.Host, mc.Port, mc.MaxPoolSize, mc.W)
}


