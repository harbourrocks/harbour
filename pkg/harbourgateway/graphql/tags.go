package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/model"
	"github.com/harbourrocks/harbour/pkg/registry/models"
)

var tagListType = graphql.NewList(tagType)
var tagType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Tag",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Description: "The name of the tag. The tag combined with the repository " +
					"is a unique identifier for an image",
			},
			"repository": &graphql.Field{
				Type:        repositoryType,
				Description: "The repository of the tag.",
			},
		},
	},
)

func TagsField(options configuration.Options) *graphql.Field {
	return &graphql.Field{
		Type:        tagListType,
		Description: "All tags of a specified repository.",
		Args: graphql.FieldConfigArgument{
			"repository": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Defines the repository of which the tags should be returned.",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// parse query parameter
			repositoryQuery, isOK := p.Args["repository"].(string)
			if !isOK {
				return nil, nil
			}

			dockerToken, err := acquireDockerToken(p.Context, options.DockerRegistry.TokenURL("repository", repositoryQuery, "pull"))

			// query tags of repository from docker registry
			var regTags models.Tags
			_, err = apiclient.Get(p.Context, options.DockerRegistry.RepositoryTagsURL(repositoryQuery), &regTags, dockerToken)
			if err != nil {
				return nil, err
			}

			// map docker registry response to harbour response
			var response = make([]model.Tag, len(regTags.Tags))
			for i, repository := range regTags.Tags {
				response[i] = model.Tag{
					Name: repository,
					Repository: model.Repository{
						Name: regTags.Name,
					},
				}
			}

			return response, err
		},
	}
}
