package handler

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourscm/models"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"net/http"
)

func Manifest(w http.ResponseWriter, r *http.Request) {
	scmConfig := configuration.GetSCMConfigReq(r)

	homepage := scmConfig.HostUrl.String()
	redirectUrl := fmt.Sprintf("%s/callback", scmConfig.HostUrl.String())
	webhookUrl := fmt.Sprintf("%s/callback", scmConfig.HostUrl.String())

	manifest := models.GithubManifest{
		Name:          "harbour.rocks",
		Url:           homepage,
		RedirectUrl:   redirectUrl,
		DefaultEvents: []models.GithubEventType{models.Push},
		HookAttributes: models.HookAttributes{
			Url:    webhookUrl,
			Active: true,
		},
		DefaultPermissions: models.DefaultPermissions{
			Contents: models.Write,
			Metadata: models.Read,
		},
	}

	_ = httphelper.WriteResponse(r, w, manifest)
}
