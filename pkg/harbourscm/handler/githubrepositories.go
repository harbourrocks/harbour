package handler

import (
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/harbourscm/github"
	"github.com/harbourrocks/harbour/pkg/harbourscm/models"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"net/http"
)

type OrganizationRepositoriesResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func OrganizationRepositories(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)
	ctx := r.Context()

	// get organization from query params
	orgLogin := httphelper.GetQueryParam(r, "org")
	if orgLogin == "" {
		log.WithField("orgLogin", orgLogin).Warn("Missing org parameter")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get installation token for organization
	token, err := github.GenerateTokenForOrganization(ctx, orgLogin)
	if err != nil {
		log.WithError(err).WithField("orgLogin", orgLogin).Error("Failed to get organization")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	repositories := make([]models.GithubRepository, 0)
	addHeaders := make(map[string]string)
	addHeaders["Authorization"] = "token " + token
	rsp, err := apiclient.Get(ctx, github.GetRepositoryUrl(orgLogin), &repositories, "", addHeaders)
	if err != nil {
		log.WithError(err).WithField("orgLogin", orgLogin).Error("Failed to get repositories")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if rsp.StatusCode >= 400 {
		log.WithField("statusCode", rsp.StatusCode).WithField("orgLogin", orgLogin).Warn("Failed to get repositories")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// map response objects
	repositoriesResponse := make([]OrganizationRepositoriesResponse, len(repositories))
	for i, repository := range repositories {
		repositoriesResponse[i] = OrganizationRepositoriesResponse{
			Id:   repository.Id,
			Name: repository.Name,
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = httphelper.WriteResponse(r, w, repositoriesResponse)
}
