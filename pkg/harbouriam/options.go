package harbouriam

import "github.com/harbourrocks/harbour/pkg/redisconfig"

// Options defines all options available to configure the IAM server.
type Options struct {
	OIDCClientID     string
	OIDCClientSecret string
	OIDCURL          string
	IAMBaseURL       string
	Redis            *redisconfig.RedisOptions
}

// NewDefaultOptions returns the default options
func NewDefaultOptions() *Options {
	s := Options{
		OIDCClientID:     "",
		OIDCClientSecret: "",
		OIDCURL:          "",
		IAMBaseURL:       "",
		Redis:            redisconfig.NewDefaultRedisOptions(),
	}

	return &s
}
