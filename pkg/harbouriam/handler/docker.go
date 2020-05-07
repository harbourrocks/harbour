package handler

import (
	"encoding/base64"
	redis2 "github.com/harbourrocks/harbour/pkg/harbouriam/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type DockerSetPassword struct {
	Password string `json:"password"`
}

// DockerModel is specific for one handler
type DockerModel struct {
	traits.HttpModel
	traits.IdTokenModel
	redisconfig.RedisModel
}

func (h DockerModel) Handle() {
	w := h.GetResponse()
	redisConfig := h.GetRedisConfig()
	idToken := h.GetToken()

	var model DockerSetPassword
	if err := h.ReadRequest(&model); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // error logged in ReadRequest
	}

	// validate password length, min 5
	if len(model.Password) < 5 {
		_ = h.WriteErrorResponse(1000)
		return
	}

	// hash the password using bcrypt, salt is automatically added during hashing
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(model.Password), 12)
	if err != nil {
		l.WithError(err).Error("Failed to hash password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode with base64 to store as string
	passwordBase64 := base64.StdEncoding.EncodeToString(passwordHashed)
	l.Tracef("Hashed password %s", passwordBase64)

	// save to redis as 'docker-password'
	client := redisconfig.OpenClient(redisConfig)
	if err := client.HSet(redis2.IamUserKey(idToken.Subject), "docker-password", passwordBase64).Err(); err != nil {
		l.WithError(err).Error("Failed to save docker-password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
