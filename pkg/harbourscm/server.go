package harbourscm

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourscm/handler"
	"github.com/harbourrocks/harbour/pkg/httppipeline"
	"github.com/sirupsen/logrus"
	"net/http"
)

// RunSCMServer runs the IAM server application
func RunSCMServer(o *configuration.Options) error {
	logrus.Info("Started Harbour SCM server")

	pipeline := httppipeline.DefaultPipeline(o.OIDCConfig, o.Redis)
	pipeline = httppipeline.WithConfig(pipeline, configuration.SCMConfigKey, *o)

	http.HandleFunc("/scm/github/manifest", pipeline(handler.Manifest))
	http.HandleFunc("/scm/github/register", pipeline(handler.GithubManualRegister))

	unPipeline := httppipeline.UnAuthPipeline(o.Redis)
	unPipeline = httppipeline.WithConfig(unPipeline, configuration.SCMConfigKey, *o)
	http.HandleFunc("/callback", unPipeline(handler.LogIncoming))
	http.HandleFunc("/scm/github/hooks", unPipeline(handler.LogIncoming))
	http.HandleFunc("/scm/github/app/redirect", unPipeline(handler.NewAppRedirect))

	http.HandleFunc("/scm/github/app", pipeline(handler.RegisterApp))

	http.HandleFunc("/scm/github/organizations", pipeline(handler.AllOrganizations))

	http.HandleFunc("/scm/github/callback", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
	})

	bindAddress := "0.0.0.0:5300"
	logrus.Info(fmt.Sprintf("Listening on httphandler://%s/", bindAddress))

	err := http.ListenAndServe(bindAddress, nil)
	logrus.WithError(err).Fatal()

	return nil
}
