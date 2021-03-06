package httppipeline

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/configuration"
	"github.com/harbourrocks/harbour/pkg/httpcontext"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"net/http"
)

func WithConfig(pipeline func(http.HandlerFunc) http.HandlerFunc, key string, config interface{}) func(http.HandlerFunc) http.HandlerFunc {
	fn := func(handler http.HandlerFunc) http.HandlerFunc {
		return pipeline(configuration.UseAddConfig(handler, key, config))
	}

	return fn
}

func DefaultPipeline(oidcConfig auth.OIDCConfig, redisConfig redisconfig.RedisOptions) func(http.HandlerFunc) http.HandlerFunc {
	fn := func(handler http.HandlerFunc) http.HandlerFunc {
		return httpcontext.
			UseRequestId(logconfig.
				UseLogger(auth.
					UseOidcTokenStr(auth.
						UseOidcToken(auth.
							UseIdToken(auth.
								UseAuth(redisconfig.
									UseRedisConfig(httpcontext.
										UseJsonResponse(handler), redisConfig))), oidcConfig))))
	}

	return fn
}

func CorsPipeline(allowedUrls []string, oidcConfig auth.OIDCConfig, redisConfig redisconfig.RedisOptions) func(http.HandlerFunc) http.HandlerFunc {
	fn := func(handler http.HandlerFunc) http.HandlerFunc {
		return httpcontext.
			UseRequestId(logconfig.
				UseLogger(httpcontext.
					UseCors(auth.
						UseOidcTokenStr(auth.
							UseOidcToken(auth.
								UseIdToken(auth.
									UseAuth(redisconfig.
										UseRedisConfig(httpcontext.
											UseJsonResponse(handler), redisConfig))), oidcConfig)), allowedUrls)))
	}

	return fn
}

func UnAuthPipeline(redisConfig redisconfig.RedisOptions) func(http.HandlerFunc) http.HandlerFunc {
	fn := func(handler http.HandlerFunc) http.HandlerFunc {
		return httpcontext.
			UseRequestId(logconfig.
				UseLogger(redisconfig.
					UseRedisConfig(httpcontext.
						UseJsonResponse(handler), redisConfig)))
	}

	return fn
}

func SemiAuthPipeline(oidcConfig auth.OIDCConfig, redisConfig redisconfig.RedisOptions) func(http.HandlerFunc) http.HandlerFunc {
	fn := func(handler http.HandlerFunc) http.HandlerFunc {
		return httpcontext.
			UseRequestId(logconfig.
				UseLogger(auth.
					UseOidcTokenStr(auth.
						UseOidcToken(auth.
							UseIdToken(redisconfig.
								UseRedisConfig(httpcontext.
									UseJsonResponse(handler), redisConfig)), oidcConfig))))
	}

	return fn
}
