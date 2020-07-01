package configuration

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/harbourrocks/harbour/pkg/registry"
	l "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Options struct {
	ContextPath    string
	Redis          redisconfig.RedisOptions
	OIDCConfig     auth.OIDCConfig
	DockerRegistry registry.RegistryConfig
	SCMConfig
}

type SCMConfig struct {
	Url string
}

func (c SCMConfig) GetCheckoutUrl() string {
	return fmt.Sprintf("%s/checkout", c.Url)
}

func NewDefaultOptions() *Options {
	s := Options{
		ContextPath: "",
		Redis:       redisconfig.NewDefaultRedisOptions(),
		OIDCConfig:  auth.DefaultConfig(),
		SCMConfig: SCMConfig{
			Url: "http://localhost:5300",
		},
	}

	return &s
}

func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	s.ContextPath = viper.GetString("CONTEXT_PATH")
	if s.ContextPath == "" {
		l.Fatal("Missing CONTEXT_PATH")
	} else if path, err := os.Stat(s.ContextPath); os.IsNotExist(err) || !path.IsDir() {
		l.Fatal("CONTEXT_PATH not found or a file")
	}

	s.SCMConfig.Url = viper.GetString("SCM_URL")
	s.SCMConfig.Url = strings.Trim(s.SCMConfig.Url, "/")

	s.OIDCConfig = auth.ParseViperConfig()
	s.Redis = redisconfig.ParseViperConfig()
	s.DockerRegistry = registry.ParseViperConfig()

	return s
}
