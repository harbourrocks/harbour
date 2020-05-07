package handler

import (
	redis2 "github.com/harbourrocks/harbour/pkg/harbouriam/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"net/http"
)

type ProfileModel struct {
	traits.HttpModel
	traits.IdTokenModel
	redisconfig.RedisModel
}

// Handle extracts the latest user information from an id token
func (h ProfileModel) Handle() {
	w := h.GetResponse()
	redisConfig := h.GetRedisConfig()
	idToken := h.GetToken()

	// save to redis as 'docker-password'
	client := redisconfig.OpenClient(redisConfig)
	err := client.HSet(redis2.IamUserKey(idToken.Subject),
		"email", idToken.Email,
		"preferred_username", idToken.PreferredUsername,
		"name", idToken.Name).Err()
	if err != nil {
		l.WithError(err).Error("Failed to save user information")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.Trace("User refreshed")

	w.WriteHeader(http.StatusOK)
}
