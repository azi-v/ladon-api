package domain

import (
	"context"

	"github.com/azi-v/ladon-api/log"
	"github.com/azi-v/ladon-api/metrics"
	"github.com/go-redis/redis/v8"
	"github.com/ory/ladon"
	"go.mongodb.org/mongo-driver/mongo"
)

var Warden *ladon.Ladon

func InitWarden(ctx context.Context,logger *log.Logger, r *redis.Client, keyPrefix string, mongo *mongo.Client) {
	// ladon需要做日志注入、缓存配置等
	cache := NewRedisManager(ctx, r, keyPrefix)
	db := NewPolicyMongoDBManager(ctx, mongo)
	Warden = &ladon.Ladon{
		Manager:     NewPolicyManager(cache, db),
		AuditLogger: logger,
		Metric:      &metrics.PrometheusMetrics{},
	}
}
