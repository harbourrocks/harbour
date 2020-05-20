package redis

import "fmt"

// GithubAppKey returns the redis key of the github app with id appId
func GithubAppKey(appId int) string {
	return fmt.Sprintf("SCM:GH:APP:%d", appId)
}

// GithubOrganizations returns the redis key of the list with all github organizations
func GithubOrganizations() string {
	return "SCM:GH:ORG:LOGIN"
}

// GithubInstallationsKey returns the redis key of the list if all app installations
func GithubInstallationsKey() string {
	return "SCM:GH:INST"
}

// GithubInstallationsKey returns the redis key of the list if all app installations
func GithubInstallationKey(installationId string) string {
	return fmt.Sprintf("SCM:GH:INST:%s", installationId)
}

// GithubOrganizationLoginKey
func GithubOrganizationLoginKey(orgLogin string) string {
	return fmt.Sprintf("SCM:GH:ORG:LOGIN:%s", orgLogin)
}

// GithubOrganizationIdKey
func GithubOrganizationIdKey(orgId string) string {
	return fmt.Sprintf("SCM:GH:ORG:ID:%s", orgId)
}
