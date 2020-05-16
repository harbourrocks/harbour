package configuration

import (
	"context"
	"github.com/harbourrocks/harbour/pkg/configuration"
	"net/http"
)

const IAMConfigKey = "iamConfig"

func GetIAMConfigReq(r *http.Request) Options {
	return GetIAMConfigCtx(r.Context())
}

func GetIAMConfigCtx(ctx context.Context) Options {
	return configuration.GetConfigCtx(ctx, IAMConfigKey).(Options)
}
