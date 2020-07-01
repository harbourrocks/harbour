package handler

import (
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"net/http"
	"strconv"
)

type RepositoryBuildsRequest struct {
	ParentKey string `json:"parent_key"`
}

type Build struct {
	BuildId     string `json:build_id`
	SCMId       string `json:"scm_id"`
	Repository  string `json:"repository"`
	Tag         string `json:"tag"`
	BuildStatus string `json:"build_status"`
	Timestamp   int64  `json:"timestamp"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	//Logs         string `json:logs`
	Commit string `json:"commit"`
}

type BuildsResponse struct {
	Builds []string `json:"builds"`
}

func RepositoryBuilds(w http.ResponseWriter, r *http.Request) {
	var repositoryBuilds RepositoryBuildsRequest
	log := logconfig.GetLogReq(r)
	client := redisconfig.GetRedisClientReq(r)

	if err := httphelper.ReadRequest(r, w, &repositoryBuilds); err != nil {
		log.WithError(err).Error("Failed to parse build request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	redisRepoEntry := client.SMembers(repositoryBuilds.ParentKey)
	if err := redisRepoEntry.Err(); err != nil {
		log.WithError(err).Error("Failed to get repo members")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	repoBuilds := redisRepoEntry.Val()

	buildsResponse := make([]Build, len(repoBuilds))
	for i, buildKey := range repoBuilds {
		buildEntry := client.HGetAll(buildKey)
		if err := buildEntry.Err(); err != nil {
			log.WithError(err).Error("Failed to get build")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		buildValue := buildEntry.Val()

		timestamp, err := strconv.ParseInt(buildValue["timestamp"], 10, 64)
		if err != nil {
			log.WithError(err).Error("Failed to parse timestamp")
		}

		startTime, err := strconv.ParseInt(buildValue["start_time"], 10, 64)
		if err != nil {
			log.WithError(err).Error("Failed to parse timestamp")
		}

		endTime, err := strconv.ParseInt(buildValue["end_time"], 10, 64)
		if err != nil {
			log.WithError(err).Error("Failed to parse timestamp")
		}

		buildsResponse[i] = Build{
			BuildId:     buildKey,
			SCMId:       buildValue["scm_id"],
			Repository:  buildValue["repository"],
			Tag:         buildValue["tag"],
			BuildStatus: buildValue["build_status"],
			Timestamp:   timestamp,
			//Logs:        buildValue["logs"],
			Commit:    buildValue["commit"],
			StartTime: startTime,
			EndTime:   endTime,
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = httphelper.WriteResponse(r, w, buildsResponse)
}
