package handler

import (
	"github.com/go-redis/redis/v7"
	"github.com/harbourrocks/harbour/pkg/auth"
	hRedis "github.com/harbourrocks/harbour/pkg/harbouriam/redis"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redis"
	"net/http"
)

// RefreshProfile extracts the latest user information from an id token
func RefreshProfile(w http.ResponseWriter, r *http.Request) {
	l := logconfig.GetLogReq(r)
	client := redisconfig.GetRedisClientReq(r)
	idToken := auth.GetIdTokenReq(r)

	var err error

	currentPrefUsername := client.HGet(hRedis.IamUserKey(idToken.Subject), "preferred_username")
	if err := currentPrefUsername.Err(); err != redis.Nil && err != nil {
		l.WithError(err).Error("Failed to load user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if no username key was found set it (this happens initially)
	userNameKey := hRedis.IamUserName(idToken.PreferredUsername)
	if currentPrefUsername.Val() == "" {
		err = client.Set(userNameKey, idToken.Subject, 0).Err()
	} else if currentPrefUsername.Val() != idToken.PreferredUsername {
		err = client.Rename(hRedis.IamUserName(currentPrefUsername.Val()), userNameKey).Err()
	}

	if err != nil {
		l.WithError(err).Error("Failed to save user information")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// save to redis as 'docker-password'
	err = client.HSet(hRedis.IamUserKey(idToken.Subject),
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
