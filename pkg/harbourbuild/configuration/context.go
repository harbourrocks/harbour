package configuration

import (
	"context"
	"github.com/harbourrocks/harbour/pkg/configuration"
	"net/http"
)

const BuildConfigKey = "buildConfig"

func GetBuildConfigReq(r *http.Request) Options {
	return GetBuildConfigCtx(r.Context())
}

func GetBuildConfigCtx(ctx context.Context) Options {
	return configuration.GetConfigCtx(ctx, BuildConfigKey).(Options)
}
