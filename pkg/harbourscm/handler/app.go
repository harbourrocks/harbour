package handler

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourscm/models"
	"github.com/harbourrocks/harbour/pkg/harbourscm/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type AppModel struct {
	traits.HttpModel
	redisconfig.RedisModel
}

func (h *AppModel) Handle(config configuration.Options) {
	w := h.GetResponse()
	r := h.GetRequest()
	redisConfig := h.GetRedisConfig()

	code := r.URL.Query().Get("code")
	if code == "" {
		l.Warningf("Invalid code from Github", code)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l.Tracef("Received code: %s", code)

	// retrieve github app configuration
	// this also returns the private key of the app which is required to build access tokens
	conversionUrl := fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code)
	appConfiguration := models.AppConfiguration{}
	_, err := apiclient.Post(conversionUrl, &appConfiguration, nil)
	if err != nil {
		// error was already logged in api client
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// save app configuration to redis
	client := redisconfig.OpenClient(redisConfig)
	err = client.HSet(redis.GithubAppKey(appConfiguration.Id),
		"clientSecret", appConfiguration.ClientSecret,
		"clientId", appConfiguration.ClientId,
		"pem", appConfiguration.PEM,
		"webhookSecret", appConfiguration.WebhookSecret,
		"name", appConfiguration.Name,
		"id", appConfiguration.Id).Err()
	if err != nil {
		l.WithError(err).Error("Failed to persist github app config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate a access token to make sure everything works
	if token, err := GenerateGithubToken(appConfiguration.Id, time.Minute*1, redisConfig); err != nil {
		l.WithError(err).Errorf("Failed to obtain access token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		l.Tracef("AccessToken: %s", token)
	}

	// setup redirect headers
	redirectUrl := fmt.Sprintf("%s/settings/vsc?app=github&id=%d", config.HostUrl, appConfiguration.Id)
	w.Header().Set("Location", redirectUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
