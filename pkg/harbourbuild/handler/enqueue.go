package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/redis"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/model"
	"github.com/harbourrocks/harbour/pkg/harbourscm/handler"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"net/http"
	"time"
)

type EnqueueHandler struct {
	config *configuration.Options
}

func NewEnqueueHandler(config *configuration.Options) EnqueueHandler {
	return EnqueueHandler{config: config}
}

func (eh EnqueueHandler) EnqueueBuild(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)
	redisClient := redisconfig.GetRedisClientReq(r)

	var buildRequest models.BuildRequest
	if err := httphelper.ReadRequest(r, w, &buildRequest); err != nil {
		log.WithError(err).Error("Failed to parse build request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buildKey, err := createBuildEntry(r.Context(), buildRequest)
	if err != nil {
		log.WithError(err).Error("Failed to save build to redis")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := handler.CheckoutRequestModel{
		SCMId:       buildRequest.SCMId,
		CallbackURL: "http://localhost:5200/build",
		State:       buildKey,
		Commit:      buildRequest.Commit,
	}

	_, err = apiclient.Post(r.Context(), "http://localhost:5300/checkout", nil, body, auth.GetOidcTokenStrCtx(r.Context()), nil)
	if err != nil {
		log.WithError(err).Error("checkout request failed")
		w.WriteHeader(http.StatusInternalServerError)
		if err := redisClient.HSet(buildKey, "build_status", "Failed").Err(); err != nil {
			log.WithError(err).Error("Failed checkout repository")
			return
		}
		return
	}

	log.Trace("Build job enqueued")
	w.WriteHeader(http.StatusOK)

	_ = httphelper.WriteResponse(r, w, model.Build{
		BuildId: buildKey,
		Status:  "Pending",
	})
}

func createBuildEntry(ctx context.Context, request models.BuildRequest) (string, error) {
	client := redisconfig.GetRedisClientCtx(ctx)

	repoKey := redis.RepoKey(request.SCMId, request.Repository)
	buildKey := redis.BuildKey(uuid.New().String())

	err := client.SAdd(repoKey, buildKey).Err()
	if err != nil {
		return buildKey, err
	}

	err = client.HSet(buildKey,
		"token", auth.GetOidcTokenStrCtx(ctx),
		"repository", request.Repository,
		"scm_id", request.SCMId,
		"commit", request.Commit,
		"logs", nil,
		"tag", request.Tag,
		"dockerfile", request.Dockerfile,
		"build_status", "Pending",
		"timestamp", time.Now().Unix()).Err()

	if err != nil {
		return buildKey, err
	}

	return buildKey, nil
}
