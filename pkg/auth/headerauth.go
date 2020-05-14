package auth

import (
	"errors"
	"github.com/coreos/go-oidc"
	l "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func HeaderAuth(request *http.Request, oidcConfig OIDCConfig) (token *oidc.IDToken, tokenString string, err error) {
	authorizationHeader := request.Header.Get("Authorization")

	// trim any spaces
	authorizationHeader = strings.TrimSpace(authorizationHeader)

	// validate header exists and is in fom "Bearer <token>"
	if authorizationHeader == "" || len(authorizationHeader) < 8 || authorizationHeader[:7] != "Bearer " {
		err = errors.New("invalid authorization header")
		l.WithField("Authorization", authorizationHeader).Error("Invalid Authorization header")
		return
	}

	l.WithField("Authorization", authorizationHeader).Trace("Authorization Header")

	// extract actual token
	tokenString = authorizationHeader[7:]

	// validate token string
	token, err = JwtAuth(tokenString, oidcConfig) // error is logged in AuthJwt

	return
}
