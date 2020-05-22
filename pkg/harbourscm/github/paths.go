package github

import "fmt"

func GetGithubApiUrl() string {
	return "https://api.github.com"
}

func GetOrganizationUrl(orgLogin string) string {
	return fmt.Sprintf("%s/orgs/%s", GetGithubApiUrl(), orgLogin)
}

func GetRepositoryUrl(orgLogin string) string {
	return fmt.Sprintf("%s/orgs/%s/repos", GetGithubApiUrl(), orgLogin)
}

func GetAppUrl() string {
	return fmt.Sprintf("%s/app", GetGithubApiUrl())
}

func GetInstallationTokenUrl(installationId string) string {
	return fmt.Sprintf("%s/app/installations/%s/access_tokens", GetGithubApiUrl(), installationId)
}
