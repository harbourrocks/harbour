package redis

import "fmt"

func BuilderAppKey(buildId int) string {
	return fmt.Sprintf("BUILD_%d", buildId)
}
