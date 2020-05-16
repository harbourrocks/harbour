package registry

import (
	"github.com/spf13/viper"
	"strings"
)

// RegistryConfig is the required minimum docker registry connections
type RegistryConfig struct {
	RegistryUrl            string
	AuthorizationServerUrl string
}

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() RegistryConfig {
	var s RegistryConfig

	s.RegistryUrl = viper.GetString("REGISTRY_URL")
	s.RegistryUrl = strings.Trim(s.RegistryUrl, " /")

	s.AuthorizationServerUrl = viper.GetString("IAM_BASE_URL")
	s.AuthorizationServerUrl = strings.Trim(s.AuthorizationServerUrl, " /")

	return s
}
