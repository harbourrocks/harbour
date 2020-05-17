package harbourbuild

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/handler"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/httppipeline"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RunBuildServer(o *configuration.Options) error {
	logrus.Info("Started Harbour build server")

	buildChan := make(chan models.BuildJob)
	builder, err := NewBuilder(buildChan, o.ContextPath, o.RepoPath, o.Redis)
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	builder.Start()

	logrus.Info("Started Harbour builder ")

	model := handler.NewBuilderModel(buildChan, o)

	pipeline := httppipeline.DefaultPipeline(o.OIDCConfig, o.Redis)
	pipeline = httppipeline.WithConfig(pipeline, configuration.BuildConfigKey, *o)

	http.HandleFunc("/build", pipeline(model.BuildImage))

	bindAddress := "127.0.0.1:5200"
	logrus.Info(fmt.Sprintf("Listening on httphandler://%s/", bindAddress))

	err = http.ListenAndServe(bindAddress, nil)
	logrus.Fatal(err)

	return err
}
