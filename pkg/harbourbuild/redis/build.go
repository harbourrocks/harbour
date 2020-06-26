package redis

import "fmt"

func BuildKey(buildId string) string {
	return fmt.Sprintf("BUILD_%s", buildId)
}

func ScmRepoKey(scmId string, repo string) string {
	return fmt.Sprintf("BUILD_%s_%s", scmId, repo)
}

func RepoKey(repo string) string {
	return fmt.Sprintf("BUILD_%s", repo)
}

func RepoTagKey(repo string, tag string) string {
	return fmt.Sprintf("BUILD_%s_%s", repo, tag)
}
