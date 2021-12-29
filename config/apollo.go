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
	CacheConfig *CacheConfig
	DBConfig    *DBConfig
}

type LogConfig struct {
}

type CacheConfig struct {
}

type MongoConfig struct {
}

type DBConfig struct {
}
