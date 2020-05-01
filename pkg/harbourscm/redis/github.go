package redis

import "fmt"

// GithubAppKey returns the redis key of the github app with id appId
func GithubAppKey(appId int) string {
	return fmt.Sprintf("SCM_GH_%d", appId)
}
