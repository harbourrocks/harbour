package handler

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"net/http"
)

type NewAppRedirectModel struct {
	OrganizationUrl string `json:"organization_url"`
}

func NewAppRedirect(w http.ResponseWriter, r *http.Request) {
	scmConfig := configuration.GetSCMConfigReq(r)
	log := logconfig.GetLogReq(r)

	var req NewAppRedirectModel
	if err := httphelper.ReadRequest(r, w, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	callbackUrl := fmt.Sprintf("%s/scm/github/fake", scmConfig.HostUrl.String())
	webhookUrl := fmt.Sprintf("%s/scm/github/hooks", scmConfig.HostUrl.String())

	githubUrl := fmt.Sprintf("%s/settings/apps/new?name=%s&url=%s&callback_url=%s&setup_url=%s&public=%s&webhook_url=%s&%s",
		req.OrganizationUrl,
		"harbour.rocks",
		scmConfig.GithubAppHomepage,
		callbackUrl,
		callbackUrl,
		"false",
		webhookUrl,
		"contents=write")

	log.WithField("url", githubUrl).Trace("Redirect url")

	w.Header().Add("Location", githubUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
