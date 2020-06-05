package configuration

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/harbourrocks/harbour/pkg/registry"
	"github.com/spf13/viper"
	"strings"
)

// Options defines all options available to configure the Gateway server.
type Options struct {
	Redis           redisconfig.RedisOptions
	OIDCConfig      auth.OIDCConfig
	DockerRegistry  registry.RegistryConfig
	CorsAllowedUrls []string
	SCMConfig
	BuildConfig
}

type SCMConfig struct {
	Url string
}

type BuildConfig struct {
	Url string
}

func (c SCMConfig) GetOrganizationsUrl() string {
	return fmt.Sprintf("%s/github/organizations", c.Url)
}

func (c SCMConfig) GetRepositoriesUrl(orgLogin string) string {
	return fmt.Sprintf("%s/github/repositories?org=%s", c.Url, orgLogin)
}

func (b BuildConfig) GetTriggerBuildUrl() string {
	return fmt.Sprintf("%s/build", b.Url)
}

// NewDefaultOptions returns the default options
func NewDefaultOptions() *Options {
	s := Options{
		Redis:      redisconfig.NewDefaultRedisOptions(),
		OIDCConfig: auth.DefaultConfig(),
		SCMConfig: SCMConfig{
			Url: "http://localhost:5300",
		},
		BuildConfig: BuildConfig{
			Url: "http://localhost:5200",
		},
	}

	return &s
}

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	s.OIDCConfig = auth.ParseViperConfig()
	s.Redis = redisconfig.ParseViperConfig()
	s.DockerRegistry = registry.ParseViperConfig()

	s.SCMConfig.Url = viper.GetString("SCM_URL")
	s.SCMConfig.Url = strings.Trim(s.SCMConfig.Url, "/")

	allowedUrls := viper.GetString("CORS_ALLOWED_URLS")
	s.CorsAllowedUrls = strings.Split(allowedUrls, ",")

	return s
}
