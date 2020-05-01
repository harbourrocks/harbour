package harbourscm

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourscm/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourscm/handler"
	"github.com/harbourrocks/harbour/pkg/httphandler"
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

		h := handler.ManifestHandler{
			Config: o,
			HttpHandler: httphandler.HttpHandler{
				Request:  r,
				Response: w,
			},
		}

		h.Handle()
	})

	http.HandleFunc("/scm/github/app", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)

		h := handler.AppHandler{
			Config: o,
			HttpHandler: httphandler.HttpHandler{
				Request:  r,
				Response: w,
			},
		}

		h.Handle()
	})

	bindAddress := "0.0.0.0:5200"
	logrus.Info(fmt.Sprintf("Listening on httphandler://%s/", bindAddress))

	err := http.ListenAndServe(bindAddress, nil)
	logrus.WithError(err).Fatal()

	return nil
}
