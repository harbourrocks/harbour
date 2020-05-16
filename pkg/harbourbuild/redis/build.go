package redis

import "fmt"

func BuildKey(buildId string) string {
	return fmt.Sprintf("BUILD_%s", buildId)
}
