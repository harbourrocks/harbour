package harbourscm

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourscm/handler"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/sirupsen/logrus"
	"net/http"
)

// RunSCMServer runs the IAM server application
func RunSCMServer(o *configuration.Options) error {
	logrus.Info("Started Harbour SCM server")

	http.HandleFunc("/scm/github/callback", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
	})

	http.HandleFunc("/scm/github/manifest", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)

		model := handler.ManifestModel{}
		traits.AddHttp(&model, r, w, o.OIDCConfig)

		model.Handle(*o)
	})

	http.HandleFunc("/scm/github/app", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)

		model := handler.AppModel{}
		traits.AddHttp(&model, r, w, o.OIDCConfig)
		redisconfig.AddRedis(&model, o.Redis)

		model.Handle(*o)
	})

	bindAddress := "0.0.0.0:5200"
	logrus.Info(fmt.Sprintf("Listening on httphandler://%s/", bindAddress))

	err := http.ListenAndServe(bindAddress, nil)
	logrus.WithError(err).Fatal()

	return nil
}
