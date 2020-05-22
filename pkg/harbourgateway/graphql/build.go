package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
)

var triggerBuildType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Build",
		Fields: graphql.Fields{
			"dockerfile": &graphql.Field{
				Type:        graphql.String,
				Description: "Name of dockerfile which should be used for build",
			},
			"tag": &graphql.Field{
				Type:        graphql.String,
				Description: "Tag which should be used for the image",
			},
			"repository": &graphql.Field{
				Type:        graphql.String,
				Description: "Code-Repo which should be built",
			},
			"commit": &graphql.Field{
				Type:        graphql.String,
				Description: "Commit which sould be used for built",
			},
		},
	})

func TriggerBuildField() *graphql.Field {
	return &graphql.Field{
		Type: triggerBuildType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			oidcTokenStr := auth.GetOidcTokenStrCtx(p.Context)

			repository, isOK := p.Args["repository"].(string)
			if !isOK {
				return nil, nil
			}

			dockerfile, isOK := p.Args["dockerfile"].(string)
			if !isOK {
				return nil, nil
			}

			tag, isOK := p.Args["tag"].(string)
			if !isOK {
				return nil, nil
			}

			commit, isOK := p.Args["commit"].(string)
			if !isOK {
				return nil, nil
			}

			build := &models.BuildRequest{
				Repository: repository,
				Dockerfile: dockerfile,
				Tag:        tag,
				Commit:     commit,
			}

			var response interface{}
			_, err := apiclient.Post(p.Context, "http://localhost:5200/build", response, build, oidcTokenStr, nil)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}
}
