package redis

import "fmt"

func BuildKey(buildId string) string {
	return fmt.Sprintf("BUILD_%s", buildId)
}

func RepoKey(scmId string, repo string) string {
	return fmt.Sprintf("BUILD_%s_%s", scmId, repo)
}
