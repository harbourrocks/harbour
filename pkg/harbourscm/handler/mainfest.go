package handler

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourscm/models"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
)

type ManifestModel struct {
	traits.HttpModel
}

func (h *ManifestModel) Handle(config configuration.Options) {
	hostUrl := config.HostUrl

	redirectUrl := fmt.Sprintf("%s/scm/github/app", hostUrl.String())
	webhookUrl := fmt.Sprintf("%s/scm/github/hooks", hostUrl.String())

	manifest := models.GithubManifest{
		Name:          "harbour.rocks",
		Url:           hostUrl.String(),
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

	_ = h.WriteResponse(manifest)
}
