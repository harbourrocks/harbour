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
	IAMConfig
}

type SCMConfig struct {
	Url string
}

type BuildConfig struct {
	Url string
}

type IAMConfig struct {
	Url string
}

func (c SCMConfig) GetOrganizationsUrl() string {
	return fmt.Sprintf("%s/github/organizations", c.Url)
}

func (c SCMConfig) GetManualRegisterUrl() string {
	return fmt.Sprintf("%s/scm/github/register", c.Url)
}

func (c SCMConfig) GetRepositoriesUrl(orgLogin string) string {
	return fmt.Sprintf("%s/github/repositories?org=%s", c.Url, orgLogin)
}

func (b BuildConfig) GetEnqueueBuildUrl() string {
	return fmt.Sprintf("%s/enqueue", b.Url)
}

func (b BuildConfig) GetRepositoryBuilds() string {
	return fmt.Sprintf("%s/builds", b.Url)
}

func (i IAMConfig) GetDockerPasswordSetUrl() string {
	return fmt.Sprintf("%s/docker/password", i.Url)
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
		IAMConfig: IAMConfig{
			Url: "http://localhost:5100",
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

	s.BuildConfig.Url = viper.GetString("BUILD_URL")
	s.BuildConfig.Url = strings.Trim(s.BuildConfig.Url, "/")

	s.IAMConfig.Url = viper.GetString("IAM_BASE_URL")
	s.IAMConfig.Url = strings.Trim(s.IAMConfig.Url, "/")

	allowedUrls := viper.GetString("CORS_ALLOWED_URLS")
	s.CorsAllowedUrls = strings.Split(allowedUrls, ",")

	return s
}
