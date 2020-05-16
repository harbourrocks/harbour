package auth

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/harbourrocks/harbour/pkg/logconfig"
)

func JwtAuth(ctx context.Context, jwtToken string, oidcConfig OIDCConfig) (token *oidc.IDToken, err error) {
	log := logconfig.GetLogCtx(ctx)

	// add token to logger only when tracing
	log.WithField("token_string", jwtToken).Trace("Parsing token")

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
		log.WithError(err).Warn("Failed to validate token")
		return
	}

	log.Info("Token valid")
	return
}
