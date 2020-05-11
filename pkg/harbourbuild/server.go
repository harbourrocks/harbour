package harbourbuild

import (
	"fmt"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/handler"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/httphandler"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RunBuildServer(o *configuration.Options) error {
	logrus.Info("Started Harbour build server")

	buildChan := make(chan models.BuildJob)
	builder, err := NewBuilder(buildChan, o.ContextPath, o.RepoPath)
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	redisconfig.AddRedis(&builder, o.Redis)
	builder.Start()
	logrus.Info("Started Harbour builder ")

	http.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
		model := handler.NewBuilderModel(buildChan)
		traits.AddHttp(&model, r, w, o.OIDCConfig)
		traits.AddIdToken(&model)
		redisconfig.AddRedis(&model, o.Redis)

		if err := httphandler.ForceAuthenticated(&model); err == nil {
			model.Handle()
		}
	})

	bindAddress := "127.0.0.1:5200"
	logrus.Info(fmt.Sprintf("Listening on httphandler://%s/", bindAddress))

	err = http.ListenAndServe(bindAddress, nil)
	logrus.Fatal(err)

	return err
}
