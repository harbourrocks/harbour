package handler

import (
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/httphandler"
	l "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type GithubRepositoryHandler struct {
	httphandler.HttpHandler
	Config *configuration.Options
}

func (h *GithubRepositoryHandler) Handle() {
	appId, _ := strconv.Atoi(h.Request.URL.Query().Get("appId"))

	// generate a access token to make sure everything works
	if token, err := GenerateGithubToken(appId, time.Minute*1, h.Config.Redis); err != nil {
		l.WithError(err).Errorf("Failed to obtain access token")
		h.Response.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		l.Tracef("AccessToken: %s", token)
	}
}
