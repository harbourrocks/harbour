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
)

type AllOrganizationsResponse struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

func AllOrganizations(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)
	client := redisconfig.GetRedisClientReq(r)
	ctx := r.Context()

	// get list with all registered github organizations
	organizationLogins, err := client.SMembers(redis.GithubOrganizations()).Result()
	if err != nil {
		log.WithError(err).Error("Failed to retrieve organization ids")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// resolve data for each organization
	allOrganizations := make([]AllOrganizationsResponse, 0)
	for _, orgLogin := range organizationLogins {
		token, err := github.GenerateTokenForOrganization(ctx, orgLogin)
		if err != nil {
			continue
		}

		org := models.GithubOrganization{}
		addHeaders := make(map[string]string)
		addHeaders["Authorization"] = "token " + token
		rsp, err := apiclient.Get(ctx, github.GetOrganizationUrl(orgLogin), &org, "", addHeaders)
		if err != nil {
			log.WithError(err).WithField("orgLogin", orgLogin).Error("Failed to get organization")
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if rsp.StatusCode >= 400 {
			log.WithField("statusCode", rsp.StatusCode).WithField("orgLogin", orgLogin).Warn("Failed to get app information")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			allOrganizations = append(allOrganizations, AllOrganizationsResponse{
				Login:     org.Login,
				Name:      org.Name,
				AvatarUrl: org.AvatarUrl,
			})
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = httphelper.WriteResponse(r, w, allOrganizations)
}
