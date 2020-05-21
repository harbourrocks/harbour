package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/redis"
	"github.com/harbourrocks/harbour/pkg/httpcontext"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	registryModels "github.com/harbourrocks/harbour/pkg/registry/models"
	"net/http"
)

type BuilderModel struct {
	buildChan chan models.BuildJob
	config    *configuration.Options
}

func NewBuilderModel(buildChan chan models.BuildJob, config *configuration.Options) BuilderModel {
	return BuilderModel{buildChan: buildChan, config: config}
}

// BuildImage
func (b BuilderModel) BuildImage(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)

	var buildRequest models.BuildRequest
	if err := httphelper.ReadRequest(r, w, &buildRequest); err != nil {
		log.WithError(err).Error("Failed to parse build request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	registryToken, err := fetchRegistryToken(r.Context(), buildRequest.Project, b.config)
	if err != nil {
		return // Error is already logged in get
	}

	buildKey, err := createBuildEntry(r.Context(), buildRequest)
	if err != nil {
		log.WithError(err).Error("Failed to save build to redis")
		return
	}

	b.buildChan <- models.BuildJob{
		Request:       buildRequest,
		BuildKey:      buildKey,
		RegistryToken: registryToken,
		RegistryUrl:   b.config.DockerRegistry.RegistryUrl,
		ReqId:         httpcontext.GetReqIdCtx(r.Context()),
	}

	log.Trace("Build job enqueued")
	w.WriteHeader(http.StatusAccepted)
}

func createBuildEntry(ctx context.Context, request models.BuildRequest) (string, error) {
	client := redisconfig.GetRedisClientCtx(ctx)

	buildId := uuid.New()
	buildKey := redis.BuildKey(buildId.String())

	err := client.HSet(buildKey,
		"build_id", buildId.String(),
		"project", request.Project,
		"commit", request.Commit,
		"logs", nil,
		"repository", request.Project,
		"build_status", "Pending").Err()

	if err != nil {
		return buildKey, err
	}

	return buildKey, nil
}

func fetchRegistryToken(ctx context.Context, repository string, registry *configuration.Options) (string, error) {
	oidcTokenStr := auth.GetOidcTokenStrCtx(ctx)
	tokenUrl := registry.DockerRegistry.TokenURL("repository", repository, "push,pull")

	var registryToken string
	var tokenResponse registryModels.DockerTokenResponse

	resp, err := apiclient.Get(ctx, tokenUrl, &tokenResponse, oidcTokenStr)
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
