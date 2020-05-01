package configuration

import (
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"net/url"
)

// Options defines all options available to configure the IAM server.
type Options struct {
	HostUrl *url.URL
	Redis   *redisconfig.RedisOptions
}

// NewDefaultOptions returns the default options
func NewDefaultOptions() *Options {
	host, err := url.Parse("httphandler://localhost:5200")

	if err != nil {
		l.WithError(err).Fatalf("Default HostUrl is invalid")
	}

	s := Options{
		HostUrl: host,
		Redis:   redisconfig.NewDefaultRedisOptions(),
	}

	return &s
}
