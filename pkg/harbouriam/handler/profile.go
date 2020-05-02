package handler

import (
	redis2 "github.com/harbourrocks/harbour/pkg/harbouriam/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"net/http"
)

type ProfileHandler struct {
	httphandler.HttpHandler
	redisconfig.RedisOptions
}

// HandleRefreshProfile extracts the latest user information from an id token
func (h ProfileHandler) HandleRefreshProfile() {
	// extract userId from token
	idToken, err := h.ExtractUser()
	if err != nil {
		h.Response.WriteHeader(http.StatusUnauthorized)
		return
	}

	// save to redis as 'docker-password'
	client := redisconfig.OpenClient(h.RedisOptions)
	err = client.HSet(redis2.IamUserKey(idToken.Subject),
		"email", idToken.Email,
		"preferred_username", idToken.PreferredUsername,
		"name", idToken.Name).Err()
	if err != nil {
		l.WithError(err).Error("Failed to save user information")
		h.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.Trace("User refreshed")

	h.Response.WriteHeader(http.StatusOK)
}
