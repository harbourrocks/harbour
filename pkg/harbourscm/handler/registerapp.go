package handler

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/harbourscm/github"
	"github.com/harbourrocks/harbour/pkg/harbourscm/models"
	"github.com/harbourrocks/harbour/pkg/harbourscm/redis"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"net/http"
	"time"
)

type RegisterAppModel struct {
	Code string `json:"code"`
}

func RegisterApp(w http.ResponseWriter, r *http.Request) {
	client := redisconfig.GetRedisClientReq(r)
	log := logconfig.GetLogReq(r)

	requestModel := RegisterAppModel{}
	if err := httphelper.ReadRequest(r, w, &requestModel); err != nil {
		log.WithError(err).Warn("Failed to read request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get code from request
	// this code will be used to get a authentication token, a webhook secret and a private key
	code := requestModel.Code
	log.WithField("code", code).Tracef("Received code")
	if code == "" {
		log.WithField("code", code).Warn("Invalid code from Github")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// retrieve github app configuration
	// this also returns the private key of the app which is required to build access tokens
	conversionUrl := fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code)
	appConfiguration := models.GithubAppConfiguration{}
	resp, err := apiclient.Post(r.Context(), conversionUrl, &appConfiguration, nil, "", nil)
	if err != nil {
		// error was already logged in api client
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if resp.StatusCode >= 400 {
		// error logged as well
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// save app configuration to redis
	err = client.HSet(redis.GithubAppKey(appConfiguration.Id),
		"pem", appConfiguration.PEM,
		"webhookSecret", appConfiguration.WebhookSecret,
		"clientSecret", appConfiguration.ClientSecret,
		"clientId", appConfiguration.ClientId,
		"updatedAt", appConfiguration.UpdatedAt,
		"createdAt", appConfiguration.CreatedAt,
		"htmlUrl", appConfiguration.HtmlUrl,
		"externalUrl", appConfiguration.ExternalUrl,
		"description", appConfiguration.Description,
		"name", appConfiguration.Name,
		"nodeId", appConfiguration.NodeId,
		"id", appConfiguration.Id,
		"owner_id", appConfiguration.Owner.ID,
		"owner_avatarUrl", appConfiguration.Owner.AvatarURL,
		"owner_login", appConfiguration.Owner.Login,
		"owner_nodeId", appConfiguration.Owner.NodeID,
		"type", appConfiguration.Owner.Type).Err()
	if err != nil {
		log.WithError(err).Error("Failed to persist github app config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// save app configuration to redis
	err = client.HSet(redis.GithubOrganizationLoginKey(appConfiguration.Owner.Login),
		"appId", appConfiguration.Id).Err()
	if err != nil {
		log.WithError(err).Error("Failed to persist github app config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = client.RPush(redis.GithubOrganizations(), appConfiguration.Owner.ID).Err()
	if err != nil {
		log.WithError(err).Error("Failed to append github organization")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate a access token to make sure everything works
	if token, err := github.GenerateGithubAppToken(r.Context(), appConfiguration.Id, appConfiguration.PEM, time.Minute*1); err != nil {
		log.WithError(err).Errorf("Failed to obtain access token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.WithField("access_token", token).Tracef("Got accessToken")
	}

	// setup redirect headers
	w.WriteHeader(http.StatusCreated)
}
