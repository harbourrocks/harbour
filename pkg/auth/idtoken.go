package auth

import (
	"github.com/coreos/go-oidc"
	l "github.com/sirupsen/logrus"
)

// IdToken harbour expects
type IdToken struct {
	// this is used for identification
	Subject string `json:"sub"`
	// candidate for docker user name
	PreferredUsername string `json:"preferred_username"`
	// candidate for docker user name
	Email string `json:"email"`
	Name  string `json:"name"`
}

// IdTokenFromToken transforms a oidc.IDToken to a harbour IDToken
func IdTokenFromToken(token *oidc.IDToken) (idToken IdToken, err error) {
	if err = token.Claims(&idToken); err != nil {
		l.WithError(err).Error("Failed to extract claims from id_token")
	}

	return
}
