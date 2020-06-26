package graphql

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/model"
)

var dockerDetailsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DockerDetails",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"passwordSet": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	})

func SetDockerPasswordField(options configuration.Options) *graphql.Field {
	return &graphql.Field{
		Type:        dockerDetailsType,
		Description: "Set the password for the docker cli",
		Args: graphql.FieldConfigArgument{
			"password": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "The password to use for cli authentication",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			oidcTokenStr := auth.GetOidcTokenStrCtx(p.Context)

			password, isOK := p.Args["password"].(string)
			if !isOK {
				return nil, errors.New("password parameter is missing")
			}

			requestModel := &model.DockerSetPasswordRequest{
				Password: password,
			}

			var response model.DockerSetPasswordResponse
			_, err := apiclient.Post(p.Context, options.IAMConfig.GetDockerPasswordSetUrl(), &response, requestModel, oidcTokenStr, nil)
			if err != nil {
				return nil, err
			}

			return response, nil
		},
	}
}
