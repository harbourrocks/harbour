package redisconfig

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"net/http"
)

const RedisConfigKey = "redisConfig"

func GetRedisConfigReq(r *http.Request) RedisOptions {
	return GetRedisConfigCtx(r.Context())
}

func GetRedisConfigCtx(ctx context.Context) RedisOptions {
	config := ctx.Value(RedisConfigKey)

	if config == nil {
		logrus.Fatal("No redis config found in context") // panics
	}

	return config.(RedisOptions)
}

func GetRedisClientReq(r *http.Request) *redis.Client {
	return GetRedisClientCtx(r.Context())
}

func GetRedisClientCtx(ctx context.Context) *redis.Client {
	redisConfig := GetRedisConfigCtx(ctx)

	client := OpenClient(redisConfig)

	return client
}

func UseRedisConfig(next http.HandlerFunc, config RedisOptions) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, RedisConfigKey, config)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}
