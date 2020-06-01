package common

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func encode(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}

func Decode(id string) string {
	bytes, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return ""
	}

	return string(bytes)
}

func GenerateGithubId(organization, repository string) string {
	return encode(fmt.Sprintf("gh/%s/%s", strings.ToLower(organization), strings.ToLower(repository)))
}

func DecomposeRepositoryId(id string) (scmProvider, second, third string) {
	id = Decode(id)

	split := strings.Split(id, "/")
	if len(split) != 3 {
		return
	}

	return split[0], split[1], split[2]
}
