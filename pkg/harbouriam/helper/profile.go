package helper

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/harbourrocks/harbour/pkg/auth"
	hRedis "github.com/harbourrocks/harbour/pkg/harbouriam/redis"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
)

type UserDetails struct {
	PreferredUsername string
	Name              string
	EMail             string
}

func RefreshProfile(ctx context.Context, idToken *auth.IdToken) (userDetails UserDetails, err error) {
	l := logconfig.GetLogCtx(ctx)
	client := redisconfig.GetRedisClientCtx(ctx)

	currentPrefUsername := client.HGet(hRedis.IamUserKey(idToken.Subject), "preferred_username")
	if err = currentPrefUsername.Err(); err != redis.Nil && err != nil {
		l.WithError(err).Error("Failed to load user")
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
		return
	}

	// save to redis as 'docker-password'
	err = client.HSet(hRedis.IamUserKey(idToken.Subject),
		"email", idToken.Email,
		"preferred_username", idToken.PreferredUsername,
		"name", idToken.Name).Err()
	if err != nil {
		l.WithError(err).Error("Failed to save user information")
		return
	}

	l.Trace("User refreshed")
	return UserDetails{
		PreferredUsername: idToken.PreferredUsername,
		Name:              idToken.Name,
		EMail:             idToken.Email,
	}, err
}
