package handler

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/harbourrocks/harbour/pkg/harbourscm/common"
	"github.com/harbourrocks/harbour/pkg/harbourscm/github"
	"github.com/harbourrocks/harbour/pkg/harbourscm/worker"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"net/http"
	"net/url"
)

type CheckoutRequestModel struct {
	SCMId       string `json:"scm_id"`
	CallbackURL string `json:"callback_url"`
	State       string `json:"state"`
	Commit      string `json:"commit"`
}

type CheckoutHandler struct {
	Github chan<- worker.GithubCheckoutTask
}

// Checkout clones the specified repository from a connected source control manager
func (h CheckoutHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)

	requestModel := CheckoutRequestModel{}
	err := httphelper.ReadRequest(r, w, &requestModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scmProvider, second, third := common.DecomposeRepositoryId(requestModel.SCMId)
	if second == "" || third == "" {
		log.Warn("Missing second and third of scmId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = url.Parse(requestModel.CallbackURL)
	if err != nil {
		log.WithField("callbackUrl", requestModel.CallbackURL).Warn("Invalid callback url")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	commitHash := plumbing.NewHash(requestModel.Commit)

	switch scmProvider {
	case "gh":
		// test github connection
		_, err := github.GenerateTokenForOrganization(r.Context(), second)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// queue task
		h.Github <- worker.GithubCheckoutTask{
			OrganizationLogin: second,
			Repository:        third,
			CallbackUrl:       requestModel.CallbackURL,
			State:             requestModel.State,
			Commit:            commitHash,
			Ctx:               r.Context(),
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
