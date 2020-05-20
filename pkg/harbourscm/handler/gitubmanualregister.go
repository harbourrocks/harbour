package handler

import (
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

type GithubManualRegisterRequest struct {
	AppId          int    `json:"app_id"`
	InstallationId string `json:"installation_id"`
	ClientId       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	PrivateKey     string `json:"private_key"`
}

const AppTokenValidity = 1 * time.Minute

// GithubManualRegister can be used to manually create or update a github app
func GithubManualRegister(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)
	client := redisconfig.GetRedisClientReq(r)

	var req GithubManualRegisterRequest
	if err := httphelper.ReadRequest(r, w, &req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // error handled in ReadRequest
	}

	appTokenStr, err := github.GenerateGithubAppToken(r.Context(), req.AppId, req.PrivateKey, AppTokenValidity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return // error logged in GenerateGithubToken
	}

	// get app configuration because we need the orgId
	appConfiguration := models.GithubAppConfiguration{}
	rsp, err := apiclient.Get(r.Context(), github.GetAppUrl(), &appConfiguration, appTokenStr, nil)
	if err != nil {
		log.WithError(err).Error("Failed to get app information")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if rsp.StatusCode >= 400 {
		log.WithField("statusCode", rsp.StatusCode).Error("Failed to get app information")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orgLogin := appConfiguration.Owner.Login

	// save app configuration to redis
	err = client.HSet(redis.GithubOrganizationLoginKey(orgLogin),
		"pem", req.PrivateKey,
		"clientSecret", req.ClientSecret,
		"installationId", req.InstallationId,
		"clientId", req.ClientId,
		"appId", req.AppId).Err()
	if err != nil {
		log.WithError(err).Error("Failed to persist github org config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// save organization to list of all organizations
	err = client.SAdd(redis.GithubOrganizations(), orgLogin).Err()
	if err != nil {
		log.WithError(err).Error("Failed to append github organization")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
