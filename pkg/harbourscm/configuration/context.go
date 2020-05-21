package configuration

import (
	"context"
	"github.com/harbourrocks/harbour/pkg/configuration"
	"net/http"
)

const SCMConfigKey = "scmConfig"

func GetSCMConfigReq(r *http.Request) Options {
	return GetSCMConfigCtx(r.Context())
}

func GetSCMConfigCtx(ctx context.Context) Options {
	return configuration.GetConfigCtx(ctx, SCMConfigKey).(Options)
}
