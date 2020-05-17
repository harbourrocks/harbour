package configuration

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/spf13/viper"
)

type Options struct {
	ContextPath string
	RepoPath    string
	Redis       redisconfig.RedisOptions
	OIDCConfig  auth.OIDCConfig
}

func NewDefaultOptions() *Options {
	s := Options{
		ContextPath: "",
		RepoPath:    "",
		Redis:       redisconfig.NewDefaultRedisOptions(),
		OIDCConfig:  auth.DefaultConfig(),
	}

	return &s
}

func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	s.ContextPath = viper.GetString("CONTEXT_PATH")
	s.RepoPath = viper.GetString("REPO_PATH")
	s.OIDCConfig = auth.ParseViperConfig()
	s.Redis = redisconfig.ParseViperConfig()
	return s
}
