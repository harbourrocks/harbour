package auth

import (
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"net/http"
	"strings"
)

func ExtractToken(request *http.Request) string {
	authorizationHeader := request.Header.Get("Authorization")
	authorizationHeader = strings.TrimSpace(authorizationHeader)

	log := logconfig.GetLogReq(request)

	// validate header exists and is in fom "Bearer <token>"
	if len(authorizationHeader) < 8 || authorizationHeader[:7] != "Bearer " {
		log.WithField("Authorization", authorizationHeader).Error("Invalid Authorization header")
		return ""
	}

	// extract actual token
	tokenString := authorizationHeader[7:]

	return tokenString
}
