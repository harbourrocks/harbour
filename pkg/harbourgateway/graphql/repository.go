package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/model"
	"github.com/harbourrocks/harbour/pkg/registry/models"
)

var repositoryListType = graphql.NewList(repositoryType)
var repositoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Repository",
		Description: "All repositories of the registry.",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Description: "The name of the docker registry repository. " +
					"This can contain slashes",
			},
		},
	},
)

func RepositoriesField(options configuration.Options) *graphql.Field {
	return &graphql.Field{
		Type: repositoryListType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// query repositories from docker registry
			var regRepositories models.Repositories
			_, err := apiclient.Get(options.DockerRegistry.RepositoriesURL(), &regRepositories)
			if err != nil {
				return nil, err
			}

			// map docker registry response to harbour response
			var response = make([]model.Repository, len(regRepositories.Repositories))
			for i, repository := range regRepositories.Repositories {
				response[i] = model.Repository{
					Name: repository,
				}
			}

			return response, err
		},
	}
}
