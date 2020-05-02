package auth

import (
	"context"
	"github.com/coreos/go-oidc"
	l "github.com/sirupsen/logrus"
)

func JwtAuth(jwtToken string, oidcConfig OIDCConfig) (token *oidc.IDToken, err error) {
	log := l.WithField("OIDC Url", oidcConfig)

	// add token to logger only when tracing
	if l.GetLevel() == l.TraceLevel {
		log = log.WithField("JWT", jwtToken)
	}

	provider, err := oidc.NewProvider(context.Background(), oidcConfig.DiscoveryUrl)
	if err != nil {
		log.WithError(err).Error("Failed to build OIDC Provider")
		return
	}

	// validate id_token
	config := oidc.Config{
		ClientID:             oidcConfig.ClientId,
		SupportedSigningAlgs: nil, // default
		SkipClientIDCheck:    false,
		SkipExpiryCheck:      false,
		SkipIssuerCheck:      false,
		Now:                  nil, // default
	}

	// validate id_token
	token, err = provider.Verifier(&config).Verify(context.Background(), jwtToken)
	if err != nil {
		log.WithError(err).Error("Failed to validate token")
		return
	}

	log.Info("Token valid")
	return
}
