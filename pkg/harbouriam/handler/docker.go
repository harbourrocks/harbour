package handler

import (
	"encoding/base64"
	"github.com/harbourrocks/harbour/pkg/harbourscm/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type DockerHandler struct {
	httphandler.HttpHandler
	redisconfig.RedisOptions
}

type DockerSetPassword struct {
	Password string `json:"password"`
}

func (h DockerHandler) HandleSetPassword() {
	var model DockerSetPassword
	if err := h.ReadRequest(&model); err != nil {
		h.Response.WriteHeader(http.StatusInternalServerError)
		return // error logged in ReadRequest
	}

	// extract userId from token
	idToken, err := h.ExtractUser()
	if err != nil {
		h.Response.WriteHeader(http.StatusUnauthorized)
		return
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
		h.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode with base64 to store as string
	passwordBase64 := base64.StdEncoding.EncodeToString(passwordHashed)
	l.Tracef("Hashed password %s", passwordBase64)

	// save to redis as 'docker-password'
	client := redisconfig.OpenClient(h.RedisOptions)
	if err := client.HSet(redis.IamUserKey(idToken.Subject), "docker-password", passwordBase64).Err(); err != nil {
		l.WithError(err).Error("Failed to save docker-password")
		h.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Response.WriteHeader(http.StatusOK)
}
