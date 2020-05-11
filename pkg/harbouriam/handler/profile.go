package handler

import (
	"github.com/go-redis/redis/v7"
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
	var err error

	client := redisconfig.OpenClient(redisConfig)

	currentPrefUsername := client.HGet(redis2.IamUserKey(idToken.Subject), "preferred_username")
	if err := currentPrefUsername.Err(); err != redis.Nil && err != nil {
		l.WithError(err).Error("Failed to load user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if no username key was found set it (this happens initially)
	userNameKey := redis2.IamUserName(idToken.PreferredUsername)
	if currentPrefUsername.Val() == "" {
		err = client.Set(userNameKey, idToken.Subject, 0).Err()
	} else if currentPrefUsername.Val() != idToken.PreferredUsername {
		err = client.Rename(redis2.IamUserName(currentPrefUsername.Val()), userNameKey).Err()
	}

	if err != nil {
		l.WithError(err).Error("Failed to save user information")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// save to redis as 'docker-password'
	err = client.HSet(redis2.IamUserKey(idToken.Subject),
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
