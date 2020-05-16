package configuration

import (
	"context"
	"github.com/harbourrocks/harbour/pkg/configuration"
	"net/http"
)

const GatewayConfigKey = "gatewayConfig"

func GetGatewayConfigReq(r *http.Request) Options {
	return GetGatewayConfigCtx(r.Context())
}

func GetGatewayConfigCtx(ctx context.Context) Options {
	return configuration.GetConfigCtx(ctx, GatewayConfigKey).(Options)
}
