package handler

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/httphandler"
	"github.com/harbourrocks/harbour/pkg/redisconfig"
	"net/http"
)

// AuthHandler handles the authentication of an user
type AuthHandler struct {
	httphandler.HttpHandler
	redisconfig.RedisOptions
}

// Test can be used to test authentication
// 401 is returned if authentication failed, 200 otherwise
func (a AuthHandler) Test() {
	_, err := auth.HeaderAuth(a.Request, a.OIDCConfig)

	if err != nil {
		a.Response.WriteHeader(http.StatusUnauthorized)
	} else {
		a.Response.WriteHeader(http.StatusOK)
	}
}
