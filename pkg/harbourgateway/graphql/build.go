package graphql

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/model"
)

var enqueueBuildType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Build",
		Fields: graphql.Fields{
			"buildId": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

func EnqueueBuildField(options configuration.Options) *graphql.Field {
	return &graphql.Field{
		Type:        enqueueBuildType,
		Description: "Enqueue build with the given information",
		Args: graphql.FieldConfigArgument{
			"repository": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Name of the Docker-Repository",
			},
			"dockerfile": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Path to dockerfile which should be used for build",
			},
			"tag": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Tag which should be used for the image",
			},
			"scmId": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "SCM-Repository which should be built",
			},
			"commit": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Commit which should be used for built",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			oidcTokenStr := auth.GetOidcTokenStrCtx(p.Context)

			repository, isOK := p.Args["repository"].(string)
			if !isOK {
				return nil, errors.New("repository parameter is missing")
			}

			dockerfile, isOK := p.Args["dockerfile"].(string)
			if !isOK {
				return nil, errors.New("dockerfile parameter is missing")
			}

			tag, isOK := p.Args["tag"].(string)
			if !isOK {
				return nil, errors.New("tag parameter is missing")
			}

			scmId, isOK := p.Args["scmId"].(string)
			if !isOK {
				return nil, errors.New("scmId parameter is missing")
			}

			commit, isOK := p.Args["commit"].(string)
			if !isOK {
				return nil, errors.New("commit parameter is missing")
			}

			build := &models.BuildRequest{
				Repository: repository,
				Dockerfile: dockerfile,
				Tag:        tag,
				SCMId:      scmId,
				Commit:     commit,
			}

			var response model.Build
			_, err := apiclient.Post(p.Context, options.BuildConfig.GetEnqueueBuildUrl(), &response, build, oidcTokenStr, nil)
			if err != nil {
				return nil, err
			}

			return response, nil
		},
	}
}
