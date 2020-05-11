package redis

import "fmt"

func BuildAppKey(buildId string) string {
	return fmt.Sprintf("BUILD_%s", buildId)
}
