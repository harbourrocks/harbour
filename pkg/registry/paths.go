package registry

import "fmt"

// RepositoriesURL returns the URL to query all repositories
//  Includes a leading slash
func RepositoriesURL() string {
	return "/v2/_catalog"
}

// RepositoriesURL returns the URL to query all repositories
func (c RegistryConfig) RepositoriesURL() string {
	return combine(c.RegistryUrl, RepositoriesURL())
}

// RepositoryTagsURL returns the URL to query all tags of a repository
//  Includes a leading slash
//  Repository is identified by its name
func RepositoryTagsURL(repositoryName string) string {
	return fmt.Sprintf("/v2/%s/tags/list", repositoryName)
}

// RepositoryTagsURL returns the URL to query all tags of a repository
//  Repository is identified by its name
func (c RegistryConfig) RepositoryTagsURL(repositoryName string) string {
	return combine(c.RegistryUrl, RepositoryTagsURL(repositoryName))
}

func TokenURL() string {
	return "/docker/auth/token"
}

func (c RegistryConfig) TokenURL(type_, name, action string) string {
	return fmt.Sprintf("%s?service=%s&scope=%s:%s:%s", combine(c.AuthorizationServerUrl, TokenURL()), c.RegistryUrl, type_, name, action)
}

func combine(host, path string) string {
	return fmt.Sprintf("%s%s", host, path)
}
