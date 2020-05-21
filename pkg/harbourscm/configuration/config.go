package configuration

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
)

// Options defines all options available to configure the IAM server.
type Options struct {
	HostUrl           *url.URL
	GithubAppHomepage *url.URL
	UiUrl             *url.URL
	Redis             redisconfig.RedisOptions
	OIDCConfig        auth.OIDCConfig
}

// NewDefaultOptions returns the default options
func NewDefaultOptions() *Options {
	host, err := url.Parse("http://localhost:5200")
	if err != nil {
		l.WithError(err).Fatalf("Default HostUrl is invalid")
	}

	githubAppHomepage, err := url.Parse("https://harbour.rocks")
	if err != nil {
		l.WithError(err).Fatalf("Default GithubAppHomepage is invalid")
	}

	uiUrl, err := url.Parse("http://localhost:4200")
	if err != nil {
		l.WithError(err).Fatalf("Default UiUrl is invalid")
	}

	s := Options{
		HostUrl:           host,
		GithubAppHomepage: githubAppHomepage,
		UiUrl:             uiUrl,
		Redis:             redisconfig.NewDefaultRedisOptions(),
		OIDCConfig:        auth.DefaultConfig(),
	}

	return &s
}

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	if u, err := url.Parse(viper.GetString("HOST_URL")); err != nil {
		l.WithError(err).Fatalf("HostUrl is invalid")
	} else {
		s.HostUrl = u
	}

	if u, err := url.Parse(viper.GetString("GITHUB_APP_HOMEPAGE")); err != nil {
		l.WithError(err).Fatalf("GithubAppHomepage is invalid")
	} else {
		s.GithubAppHomepage = u
	}

	if u, err := url.Parse(viper.GetString("UI_URL")); err != nil {
		l.WithError(err).Fatalf("UiUrl is invalid")
	} else {
		s.UiUrl = u
	}

	s.Redis = redisconfig.ParseViperConfig()
	s.OIDCConfig = auth.ParseViperConfig()

	return s
}
