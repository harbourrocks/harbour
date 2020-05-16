package handler

import (
	"github.com/harbourrocks/harbour/pkg/httpcontext/traits"
	"github.com/harbourrocks/harbour/pkg/redis"
	l "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type GithubRepositoryModel struct {
	traits.HttpModel
	redisconfig.RedisModel
}

func (h *GithubRepositoryModel) Handle() {
	r := h.GetRequest()
	w := h.GetResponse()
	redisConfig := h.GetRedisConfig()

	appId, _ := strconv.Atoi(r.URL.Query().Get("appId"))

	// generate a access token to make sure everything works
	if token, err := GenerateGithubToken(appId, time.Minute*1, redisConfig); err != nil {
		l.WithError(err).Errorf("Failed to obtain access token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		l.Tracef("AccessToken: %s", token)
	}
}
