package registry

import (
	"github.com/spf13/viper"
	"strings"
)

// RegistryConfig is the required minimum docker registry connections
type RegistryConfig struct {
	Url string
}

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() RegistryConfig {
	var s RegistryConfig

	s.Url = viper.GetString("REGISTRY_URL")
	s.Url = strings.Trim(s.Url, " /")

	return s
}
