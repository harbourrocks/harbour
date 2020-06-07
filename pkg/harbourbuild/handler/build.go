package handler

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/harbourscm/worker"
	"github.com/harbourrocks/harbour/pkg/httpcontext"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	registryModels "github.com/harbourrocks/harbour/pkg/registry/models"
	"net/http"
)

type BuildHandler struct {
	buildChan chan models.BuildJob
	config    *configuration.Options
}

func NewBuildHandler(buildChan chan models.BuildJob, config *configuration.Options) BuildHandler {
	return BuildHandler{buildChan: buildChan, config: config}
}

func (b BuildHandler) Build(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)
	client := redisconfig.GetRedisClientReq(r)

	buildKey := httphelper.GetQueryParam(r, "state")

	var checkoutResponse worker.CheckoutCompletedModel
	if err := httphelper.ReadRequest(r, w, &checkoutResponse); err != nil {
		log.WithError(err).Error("Failed to parse build request")
		if err := client.HSet(buildKey, "build_status", "Failed").Err(); err != nil {
			log.WithError(err).Error("Failed to save build data")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if checkoutResponse.Success != true {
		if err := client.HSet(buildKey, "build_status", "Failed").Err(); err != nil {
			log.WithError(err).Error("Failed to save build data")
			return
		}
	}

	redisBuildEntry := client.HGetAll(buildKey)
	if err := redisBuildEntry.Err(); err != redis.Nil && err != nil {
		if err := client.HSet(buildKey, "build_status", "Failed").Err(); err != nil {
			log.WithError(err).Error("Failed to save build data")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	buildEntry := redisBuildEntry.Val()

	ctx := context.WithValue(r.Context(), "oidcTokenStr", buildEntry["token"])

	repository := buildEntry["repository"]

	registryToken, err := fetchRegistryToken(ctx, repository, b.config)
	if err != nil {
		if err := client.HSet(buildKey, "build_status", "Failed").Err(); err != nil {
			log.WithError(err).Error("Failed to save build data")
			return
		}
		return // Error is already logged in get
	}

	b.buildChan <- models.BuildJob{
		Repository:    repository,
		Tag:           buildEntry["tag"],
		FilePath:      checkoutResponse.WorkspacePath,
		Dockerfile:    buildEntry["dockerfile"],
		BuildKey:      buildKey,
		RegistryToken: registryToken,
		RegistryUrl:   b.config.DockerRegistry.RegistryUrl,
		ReqId:         httpcontext.GetReqIdCtx(r.Context()),
	}

	log.Trace("Build job enqueued")
	w.WriteHeader(http.StatusOK)
}

func fetchRegistryToken(ctx context.Context, repository string, registry *configuration.Options) (string, error) {
	oidcTokenStr := auth.GetOidcTokenStrCtx(ctx)
	tokenUrl := registry.DockerRegistry.TokenURL("repository", repository, "push,pull")

	var registryToken string
	var tokenResponse registryModels.DockerTokenResponse

	resp, err := apiclient.Get(ctx, tokenUrl, &tokenResponse, oidcTokenStr, nil)
	if err != nil {
		return registryToken, err
	}

	if resp.StatusCode >= 400 {
		err = errors.New("request failed")
		return registryToken, err
	}

	registryToken = tokenResponse.Token
	return registryToken, nil
}
