package configuration

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetConfigReq(r *http.Request, key string) interface{} {
	return GetConfigCtx(r.Context(), key)
}

func GetConfigCtx(ctx context.Context, key string) interface{} {
	config := ctx.Value(key)

	if config == nil {
		logrus.WithField("key", key).Fatal("No config found in context") // panics
	}

	return config
}

func UseAddConfig(next http.HandlerFunc, key string, config interface{}) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, key, config)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}
