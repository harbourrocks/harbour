package graphql

import (
	"bytes"
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourbuild/models"
)

var triggerBuildType = graphql.NewObject(
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

func TriggerBuildField() *graphql.Field {
	return &graphql.Field{
		Type:        triggerBuildType,
		Description: "Trigger build with the given information",
		Args: graphql.FieldConfigArgument{
			"dockerfile": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Name of dockerfile which should be used for build",
			},
			"tag": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Tag which should be used for the image",
			},
			"repository": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Code-Repo which should be built",
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

			body := new(bytes.Buffer)
			err := json.NewEncoder(body).Encode(build)
			if err != nil {
				return nil, err
			}

			var response interface{}
			_, err = apiclient.Post(p.Context, "http://localhost:5200/build", response, body, oidcTokenStr, nil)
			if err != nil {
				return nil, err
			}

			return response, nil
		},
	}
}
