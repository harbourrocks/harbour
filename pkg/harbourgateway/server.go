package server

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/graphqlcontext"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/configuration"
	graphql2 "github.com/harbourrocks/harbour/pkg/harbourgateway/graphql"
	"github.com/harbourrocks/harbour/pkg/httppipeline"
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
				"repositories":        graphql2.RepositoriesField(*o),
				"tags":                graphql2.TagsField(*o),
				"githubOrganizations": graphql2.GithubOrganizationsField(*o),
				"githubRepositories":  graphql2.GithubRepositoriesField(*o),
			},
		})

	var mutationType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"enqueueBuild": graphql2.EnqueueBuildField(*o),
			},
		})

	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	pipeline := httppipeline.CorsPipeline(o.CorsAllowedUrls, o.OIDCConfig, o.Redis)
	pipeline = httppipeline.WithConfig(pipeline, configuration.GatewayConfigKey, *o)

	http.HandleFunc("/graphql", pipeline(graphqlcontext.UseGraphQl(schema)))

	bindAddress := "127.0.0.1:5400"
	logrus.Info(fmt.Sprintf("Listening on http://%s/", bindAddress))

	err := http.ListenAndServe(bindAddress, nil)
	logrus.Fatal(err)

	return err
}
