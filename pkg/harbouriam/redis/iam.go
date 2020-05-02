package redis

import "fmt"

// IamUserKey returns the redis key of the specified user
func IamUserKey(userId string) string {
	return fmt.Sprintf("IAM_USER_%s", userId)
}
