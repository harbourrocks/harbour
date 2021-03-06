package graphql

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/apiclient"
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/harbourgateway/configuration"
	"github.com/harbourrocks/harbour/pkg/harbourscm/handler"
)

var githubOrganizationListType = graphql.NewList(githubOrganizationType)
var githubOrganizationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "GithubOrganization",
		Description: "A GitHub organization liked to harbour.",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type:        graphql.String,
				Description: "The identification for a github organization.",
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The display name of the organization.",
			},
			"avatarUrl": &graphql.Field{
				Type:        graphql.String,
				Description: "A url to the image of the avatar.",
			},
		},
	},
)

type SCMGithubOrganizationsResponse struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

func GithubOrganizationsField(options configuration.Options) *graphql.Field {
	return &graphql.Field{
		Type: githubOrganizationListType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			oidcToken := auth.GetOidcTokenStrCtx(p.Context)

			// query repositories from docker registry
			organizationsResponse := make([]SCMGithubOrganizationsResponse, 0)
			rsp, err := apiclient.Get(p.Context, options.SCMConfig.GetOrganizationsUrl(), &organizationsResponse, oidcToken, nil)
			if err != nil || rsp.StatusCode >= 300 {
				return nil, err
			}

			return organizationsResponse, err
		},
	}
}

var githubRepositoryListType = graphql.NewList(githubRepositoryType)
var githubRepositoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "GithubRepository",
		Description: "A github repository of a specific organisation.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "The identification for a github repository.",
			},
			"scm_id": &graphql.Field{
				Type:        graphql.String,
				Description: "A unique identifier withing whole harbour.",
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The display name of the repository.",
			},
		},
	},
)

type SCMGithubRepositoriesResponse struct {
	Id    int    `json:"id"`
	SCMId string `json:"scm_id"`
	Name  string `json:"name"`
}

func GithubRepositoriesField(options configuration.Options) *graphql.Field {
	return &graphql.Field{
		Type: githubRepositoryListType,
		Args: graphql.FieldConfigArgument{
			"orgLogin": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Defines the organization which repositories to return.",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			oidcToken := auth.GetOidcTokenStrCtx(p.Context)

			// parse query parameter
			login, isOK := p.Args["orgLogin"].(string)
			if !isOK {
				return nil, errors.New("login parameter is missing")
			}

			// query repositories from docker registry
			organizationsResponse := make([]SCMGithubRepositoriesResponse, 0)
			rsp, err := apiclient.Get(p.Context, options.SCMConfig.GetRepositoriesUrl(login), &organizationsResponse, oidcToken, nil)
			if err != nil || rsp.StatusCode >= 300 {
				return nil, err
			}

			return organizationsResponse, err
		},
	}
}

var githubRegisterAppType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "RegisterApp",
		Description: "Type to register an app manually.",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type:        graphql.Int,
				Description: "The Status-Code for the request",
			},
		},
	},
)

type RegisterAppResponse struct {
	Status int `json:status`
}

func GithubRegisterAppField(options configuration.Options) *graphql.Field {
	return &graphql.Field{
		Type: githubRegisterAppType,
		Args: graphql.FieldConfigArgument{
			"appId": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "The Id for the github application.",
			},
			"installationId": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "The installation_id of the application.",
			},
			"clientId": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "The client_id of the application.",
			},
			"clientSecret": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "The client_secret of the application.",
			},
			"privateKey": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "The private_key of the application",
			},
		}, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			oidcToken := auth.GetOidcTokenStrCtx(p.Context)

			appId, isOk := p.Args["appId"].(int)
			if !isOk {
				return nil, errors.New("appId parameter is missing")
			}

			clientId, isOk := p.Args["clientId"].(string)
			if !isOk {
				return nil, errors.New("clientId parameter is missing")
			}

			installationId, isOk := p.Args["installationId"].(string)
			if !isOk {
				return nil, errors.New("installationId parameter is missing")
			}

			clientSecret, isOk := p.Args["clientSecret"].(string)
			if !isOk {
				return nil, errors.New("clientSecret parameter is missing")
			}

			privateKey, isOk := p.Args["privateKey"].(string)
			if !isOk {
				return nil, errors.New("privateKey parameter is missing")
			}

			request := &handler.GithubManualRegisterRequest{
				AppId:          appId,
				InstallationId: installationId,
				ClientId:       clientId,
				ClientSecret:   clientSecret,
				PrivateKey:     privateKey,
			}

			rsp, err := apiclient.Post(p.Context, options.SCMConfig.GetManualRegisterUrl(), nil, request, oidcToken, nil)
			if err != nil || rsp.StatusCode >= 300 {
				return RegisterAppResponse{Status: rsp.StatusCode}, err
			}
			return RegisterAppResponse{Status: rsp.StatusCode}, err
		},
	}
}
