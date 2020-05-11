package handler

import (
	"encoding/json"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
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
	//redisConfig := b.GetRedisConfig()
	var buildRequest models.BuildRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&buildRequest)
	if err != nil {
		l.WithError(err).Error("Failed to parse build request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b.buildChan <- models.BuildJob{Request: buildRequest}

	l.Trace("Build job enqueued")

	w.WriteHeader(http.StatusAccepted)
}
