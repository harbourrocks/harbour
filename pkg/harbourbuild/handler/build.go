package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/redis"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redis"
	"net/http"
)

type BuilderModel struct {
	buildChan chan models.BuildJob
}

func NewBuilderModel(buildChan chan models.BuildJob) BuilderModel {
	return BuilderModel{buildChan: buildChan}
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

	buildKey, err := createBuildEntry(r.Context(), buildRequest)
	if err != nil {
		log.WithError(err).Error("Failed to save build to redis")
		return
	}

	b.buildChan <- models.BuildJob{Request: buildRequest, BuildKey: buildKey}

	log.Trace("Build job enqueued")
	w.WriteHeader(http.StatusAccepted)
}

func createBuildEntry(ctx context.Context, request models.BuildRequest) (string, error) {
	client := redisconfig.GetRedisClientCtx(ctx)

	buildId := uuid.New()
	buildKey := redis.BuildAppKey(buildId.String())

	err := client.HSet(buildKey,
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
