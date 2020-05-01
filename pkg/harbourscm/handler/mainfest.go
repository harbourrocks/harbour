package handler

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/httphandler"
)

type ManifestHandler struct {
	httphandler.HttpHandler
	Config *configuration.Options
}

func (h *ManifestHandler) Handle() {
	hostUrl := h.Config.HostUrl

	redirectUrl := fmt.Sprintf("%s/scm/github/app", hostUrl.String())
	webhookUrl := fmt.Sprintf("%s/scm/github/hooks", hostUrl.String())

	manifest := GithubManifest{
		Name:          "harbour.rocks",
		Url:           hostUrl.String(),
		RedirectUrl:   redirectUrl,
		DefaultEvents: []GithubEventType{Push},
		HookAttributes: HookAttributes{
			Url:    webhookUrl,
			Active: true,
		},
		DefaultPermissions: DefaultPermissions{
			Contents: Write,
			Metadata: Read,
		},
	}

	_ = h.WriteResponse(manifest)
}
