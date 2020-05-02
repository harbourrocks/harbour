package harbouriam

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbouriam/configuration"
	"github.com/harbourrocks/harbour/pkg/harbouriam/handler"
	"github.com/harbourrocks/harbour/pkg/httphandler"
	"net/http"
	"net/url"
	"path"

	"github.com/sirupsen/logrus"
)

// RunIAMServer runns the IAM server application
func RunIAMServer(o *configuration.Options) error {
	logrus.Info("Started Harbour IAM server")

	// obtain login redirect url
	redirectURL, err := url.Parse(o.IAMBaseURL)
	if err != nil {
		logrus.Fatal(err)
	} else {
		redirectURL.Path = path.Join(redirectURL.Path, "/auth/oidc/callback")
	}

	http.HandleFunc("/auth/test", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
		// AuthHandler
		authHandler := handler.AuthHandler{
			HttpHandler: httphandler.HttpHandler{
				Request:    r,
				Response:   w,
				OIDCConfig: o.OIDCConfig,
			},
			RedisOptions: o.Redis,
		}
		authHandler.Test()
	})

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
		// AuthHandler
		profileHandler := handler.ProfileHandler{
			HttpHandler: httphandler.HttpHandler{
				Request:    r,
				Response:   w,
				OIDCConfig: o.OIDCConfig,
			},
			RedisOptions: o.Redis,
		}
		profileHandler.HandleRefreshProfile()
	})

	// DockerHandler

	http.HandleFunc("/docker/password", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
		dockerHandler := handler.DockerHandler{
			HttpHandler: httphandler.HttpHandler{
				Request:    r,
				Response:   w,
				OIDCConfig: o.OIDCConfig,
			},
			RedisOptions: o.Redis,
		}
		dockerHandler.HandleSetPassword()
	})

	bindAddress := "127.0.0.1:5100"
	logrus.Info(fmt.Sprintf("Listening on httphandler://%s/", bindAddress))

	err = http.ListenAndServe(bindAddress, nil)
	logrus.Fatal(err)

	return err
}
