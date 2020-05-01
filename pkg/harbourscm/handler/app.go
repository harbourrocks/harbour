package handler

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourscm/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type AppHandler struct {
	httphandler.HttpHandler
	Config *configuration.Options
}

func (h *AppHandler) Handle() {
	code := h.Request.URL.Query().Get("code")
	if code == "" {
		l.Warningf("Invalid code from Github", code)
		h.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Tracef("Received code: %s", code)

	// retrieve github app configuration
	// this also returns the private key of the app which is required to build access tokens
	conversionUrl := fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code)
	appConfiguration := AppConfiguration{}
	_, err := apiclient.Post(conversionUrl, &appConfiguration, nil)
	if err != nil {
		// error was already logged in api client
		h.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// save app configuration to redis
	client := redisconfig.OpenClient(h.Config.Redis)
	err = client.HSet(redis.GithubAppKey(appConfiguration.Id),
		"clientSecret", appConfiguration.ClientSecret,
		"clientId", appConfiguration.ClientId,
		"pem", appConfiguration.PEM,
		"webhookSecret", appConfiguration.WebhookSecret,
		"name", appConfiguration.Name,
		"id", appConfiguration.Id).Err()
	if err != nil {
		l.WithError(err).Error("Failed to persist github app config")
		h.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate a access token to make sure everything works
	if token, err := GenerateGithubToken(appConfiguration.Id, time.Minute*1, h.Config.Redis); err != nil {
		l.WithError(err).Errorf("Failed to obtain access token")
		h.Response.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		l.Tracef("AccessToken: %s", token)
	}

	// setup redirect headers
	redirectUrl := fmt.Sprintf("%s/settings/vsc?app=github&id=%d", h.Config.HostUrl, appConfiguration.Id)
	h.Response.Header().Set("Location", redirectUrl)
	h.Response.WriteHeader(http.StatusTemporaryRedirect)
}
