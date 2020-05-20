package graphql

import (
	"context"
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
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

func acquireDockerToken(ctx context.Context, tokenUrl string) (token string, err error) {
	oidcTokenStr := auth.GetOidcTokenStrCtx(ctx)

	var tokenResponse models.DockerTokenResponse
	resp, err := apiclient.Get(ctx, tokenUrl, &tokenResponse, oidcTokenStr, nil)
	if err != nil {
		return // error logged in Get
	}

	if resp.StatusCode >= 400 {
		err = errors.New("request failed")
		return // error logged in
	}

	token = tokenResponse.Token
	return
}

func RepositoriesField(options configuration.Options) *graphql.Field {
	return &graphql.Field{
		Type: repositoryListType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			dockerToken, err := acquireDockerToken(p.Context, options.DockerRegistry.TokenURL("registry", "catalog", "*"))

			// query repositories from docker registry
			var regRepositories models.Repositories
			rsp, err := apiclient.Get(p.Context, options.DockerRegistry.RepositoriesURL(), &regRepositories, dockerToken)
			if err != nil || rsp.StatusCode >= 300 {
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
