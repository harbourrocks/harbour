package harbouriam

import "github.com/spf13/viper"

// ParseViperConfig tries to map a viper configuration
func ParseViperConfig() *Options {
	s := NewDefaultOptions()

	s.OIDCClientID = viper.GetString("OIDC_CLIENT_ID")
	s.OIDCClientSecret = viper.GetString("OIDC_CLIENT_SECRET")
	s.OIDCURL = viper.GetString("OIDC_URL")
	s.IAMBaseURL = viper.GetString("IAM_BASE_URL")

	return s
}
