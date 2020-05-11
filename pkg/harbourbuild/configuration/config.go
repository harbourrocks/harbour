package configuration

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
)

type Options struct {
	Redis      redisconfig.RedisOptions
	OIDCConfig auth.OIDCConfig
}

func NewDefaultOptions() *Options {
	s := Options{
		Redis:      redisconfig.NewDefaultRedisOptions(),
		OIDCConfig: auth.DefaultConfig(),
	}

	return &s
}

func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	s.OIDCConfig = auth.ParseViperConfig()
	s.Redis = redisconfig.ParseViperConfig()
	return s
}
