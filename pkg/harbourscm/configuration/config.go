package configuration

import (
	l "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
)

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	if hostUrl, err := url.Parse(viper.GetString("HOST_URL")); err != nil {
		l.WithError(err).Fatalf("HostUrl is invalid")
	} else {
		s.HostUrl = hostUrl
	}

	return s
}
