package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/redis"
	"github.com/harbourrocks/harbour/pkg/harbourscm/handler"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"net/http"
)

type EnqueuHandler struct {
	config *configuration.Options
}

func NewEnqueueHandler(config *configuration.Options) EnqueuHandler {
	return EnqueuHandler{config: config}
}

func (eh EnqueuHandler) EnqueueBuild(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)

	var buildRequest models.BuildRequest
	if err := httphelper.ReadRequest(r, w, &buildRequest); err != nil {
		log.WithError(err).Error("Failed to parse build request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buildKey, err := createBuildEntry2(r.Context(), buildRequest)
	if err != nil {
		log.WithError(err).Error("Failed to save build to redis")
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
		return
	}

	log.Trace("Build enqueued")
	w.WriteHeader(http.StatusOK)
}

func createBuildEntry2(ctx context.Context, request models.BuildRequest) (string, error) {
	client := redisconfig.GetRedisClientCtx(ctx)

	buildId := uuid.New()
	buildKey := redis.BuildKey(buildId.String())

	err := client.HSet(buildKey,
		"build_id", buildId.String(),
		"repository", request.SCMId,
		"commit", request.Commit,
		"logs", nil,
		"build_status", "Pending").Err()

	if err != nil {
		return buildKey, err
	}

	return buildKey, nil
}
