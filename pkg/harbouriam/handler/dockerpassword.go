package handler

import (
	"encoding/base64"
	"github.com/harbourrocks/harbour/pkg/auth"
	hRedis "github.com/harbourrocks/harbour/pkg/harbouriam/redis"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redis"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type DockerSetPassword struct {
	Password string `json:"password"`
}

// DockerPassword
func DockerPassword(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)
	client := redisconfig.GetRedisClientReq(r)
	idToken := auth.GetIdTokenReq(r)

	var model DockerSetPassword
	if err := httphelper.ReadRequest(r, w, &model); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // error logged in ReadRequest
	}

	// validate password length, min 5
	if len(model.Password) < 5 {
		_ = httphelper.WriteErrorResponse(r, w, 1000)
		return
	}

	// hash the password using bcrypt, salt is automatically added during hashing
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(model.Password), 12)
	if err != nil {
		log.WithError(err).Error("Failed to hash password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// encode with base64 to store as string
	passwordBase64 := base64.StdEncoding.EncodeToString(passwordHashed)
	log.Tracef("Hashed password %s", passwordBase64)

	// save to redis as 'docker-password'
	if err := client.HSet(hRedis.IamUserKey(idToken.Subject), "docker-password", passwordBase64).Err(); err != nil {
		log.WithError(err).Error("Failed to save docker-password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
