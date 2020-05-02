package auth

import "github.com/spf13/viper"

// OIDCConfig is the required minimum for token validation
type OIDCConfig struct {
	DiscoveryUrl string
	ClientId     string
}

func DefaultConfig() OIDCConfig {
	return OIDCConfig{
		DiscoveryUrl: "",
		ClientId:     "",
	}
}

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() OIDCConfig {
	var s OIDCConfig

	s.ClientId = viper.GetString("OIDC_CLIENT_ID")
	s.DiscoveryUrl = viper.GetString("OIDC_DISCOVERY_URL")

	return s
}
