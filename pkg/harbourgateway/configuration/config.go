package configuration

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/redis"
	"github.com/harbourrocks/harbour/pkg/registry"
)

// Options defines all options available to configure the Gateway server.
type Options struct {
	Redis          redisconfig.RedisOptions
	OIDCConfig     auth.OIDCConfig
	DockerRegistry registry.RegistryConfig
}

// NewDefaultOptions returns the default options
func NewDefaultOptions() *Options {
	s := Options{
		Redis:      redisconfig.NewDefaultRedisOptions(),
		OIDCConfig: auth.DefaultConfig(),
	}

	return &s
}

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	s.OIDCConfig = auth.ParseViperConfig()
	s.Redis = redisconfig.ParseViperConfig()
	s.DockerRegistry = registry.ParseViperConfig()

	return s
}
