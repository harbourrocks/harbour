package harbouriam

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbouriam/configuration"
	"github.com/harbourrocks/harbour/pkg/harbouriam/handler"
	"github.com/harbourrocks/harbour/pkg/httppipeline"
	"github.com/sirupsen/logrus"
	"net/http"
)

// RunIAMServer runs the IAM server application
func RunIAMServer(o *configuration.Options) error {
	logrus.Info("Started Harbour IAM server")

	pipeline := httppipeline.DefaultPipeline(o.OIDCConfig, o.Redis)
	pipeline = httppipeline.WithConfig(pipeline, configuration.IAMConfigKey, *o)

	http.HandleFunc("/refresh", pipeline(handler.RefreshProfile))

	http.HandleFunc("/docker/password", pipeline(handler.DockerPassword))

	unAuthPipeline := httppipeline.UnAuthPipeline(o.Redis)
	unAuthPipeline = httppipeline.WithConfig(unAuthPipeline, configuration.IAMConfigKey, *o)
	http.HandleFunc("/docker/auth/token", unAuthPipeline(handler.DockerToken))

	bindAddress := "0.0.0.0:5100"
	logrus.Info(fmt.Sprintf("Listening on http://%s/", bindAddress))

	err := http.ListenAndServe(bindAddress, nil)
	logrus.Fatal(err)

	return err
}
