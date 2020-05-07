package server

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/configuration"
	graphql2 "github.com/harbourrocks/harbour/pkg/harbourgateway/graphql"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/handler"
	traits2 "github.com/harbourrocks/harbour/pkg/harbourgateway/traits"
	"github.com/harbourrocks/harbour/pkg/httphandler"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"github.com/sirupsen/logrus"
	"net/http"
)

// RunGatewayServer runs the Gateway server application
func RunGatewayServer(o *configuration.Options) error {
	logrus.Info("Started Harbour Gateway server")

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"repositories": graphql2.RepositoriesField(*o),
				"tags":         graphql2.TagsField(*o),
			},
		})

	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		logrus.Trace(r)
		model := handler.GraphQLModel{}
		traits.AddHttp(&model, r, w, o.OIDCConfig)
		traits.AddIdToken(&model)
		traits2.AddGraphQL(&model, schema)
		if err := httphandler.ForceAuthenticated(&model); err != nil {
			_ = model.Handle()
		}
	})

	bindAddress := "127.0.0.1:5400"
	logrus.Info(fmt.Sprintf("Listening on http://%s/", bindAddress))

	err := http.ListenAndServe(bindAddress, nil)
	logrus.Fatal(err)

	return err
}
