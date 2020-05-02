package configuration

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/spf13/viper"
)

// Options defines all options available to configure the IAM server.
type Options struct {
	OIDCClientID     string
	OIDCClientSecret string
	OIDCURL          string
	IAMBaseURL       string
	Redis            redisconfig.RedisOptions
	OIDCConfig       auth.OIDCConfig
}

// NewDefaultOptions returns the default options
func NewDefaultOptions() *Options {
	s := Options{
		OIDCClientID:     "",
		OIDCClientSecret: "",
		OIDCURL:          "",
		IAMBaseURL:       "",
		Redis:            redisconfig.NewDefaultRedisOptions(),
		OIDCConfig:       auth.DefaultConfig(),
	}

	return &s
}

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	s.OIDCClientID = viper.GetString("OIDC_CLIENT_ID")
	s.OIDCClientSecret = viper.GetString("OIDC_CLIENT_SECRET")
	s.OIDCURL = viper.GetString("OIDC_URL")
	s.IAMBaseURL = viper.GetString("IAM_BASE_URL")

	s.OIDCConfig = auth.ParseViperConfig()
	s.Redis = redisconfig.ParseViperConfig()

	return s
}
