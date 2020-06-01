package configuration

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/harbourrocks/harbour/pkg/registry"
	"github.com/spf13/viper"
)

type Options struct {
	ContextPath    string
	Redis          redisconfig.RedisOptions
	OIDCConfig     auth.OIDCConfig
	DockerRegistry registry.RegistryConfig
}

func NewDefaultOptions() *Options {
	s := Options{
		ContextPath: "",
		Redis:       redisconfig.NewDefaultRedisOptions(),
		OIDCConfig:  auth.DefaultConfig(),
	}

	return &s
}

func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	s.ContextPath = viper.GetString("CONTEXT_PATH")
	s.OIDCConfig = auth.ParseViperConfig()
	s.Redis = redisconfig.ParseViperConfig()
	s.DockerRegistry = registry.ParseViperConfig()

	return s
}
