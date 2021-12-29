package domain

import (
	"context"
	"time"

	"git.ymt360.com/usercenter/ymt-ladon/log"
	"git.ymt360.com/usercenter/ymt-ladon/metrics"
	"github.com/go-redis/redis/v8"
	"github.com/ory/ladon"
)

var Warden *ladon.Ladon

func Init(logger *log.Logger, r *redis.Client, keyPrefix string) {
	// ladon需要做日志注入、缓存配置等
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cache := NewRedisManager(ctx, r, keyPrefix)
	db := NewPolicyDBManager()
	Warden = &ladon.Ladon{
		Manager:     NewPolicyManager(cache, db),
		AuditLogger: logger,
		Metric: &metrics.PrometheusMetrics{},
	}
}
