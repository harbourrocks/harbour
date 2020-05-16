package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/redis"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	l "github.com/sirupsen/logrus"
	"net/http"
)

type BuilderModel struct {
	traits.HttpModel
	traits.IdTokenModel
	redisconfig.RedisModel
	buildChan chan models.BuildJob
}

func NewBuilderModel(buildChan chan models.BuildJob) BuilderModel {
	return BuilderModel{buildChan: buildChan}
}

func (b BuilderModel) Handle() {
	w := b.GetResponse()
	req := b.GetRequest()
	var buildRequest models.BuildRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&buildRequest)
	if err != nil {
		l.WithError(err).Error("Failed to parse build request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buildKey, err := b.createBuildEntry(buildRequest)
	if err != nil {
		l.WithError(err).Error("Failed to save build to redis")
		return
	}

	b.buildChan <- models.BuildJob{Request: buildRequest, BuildKey: buildKey}

	l.Trace("Build job enqueued")
	w.WriteHeader(http.StatusAccepted)
}

func (b BuilderModel) createBuildEntry(request models.BuildRequest) (string, error) {
	redisConfig := b.GetRedisConfig()
	buildId := uuid.New()
	buildKey := redis.BuildKey(buildId.String())

	client := redisconfig.OpenClient(redisConfig)
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
