package handler

import (
	"github.com/harbourrocks/harbour/pkg/auth"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"net/http"
)

// AuthModel handles the authentication of an user
type AuthModel struct {
	traits.HttpModel
	traits.IdTokenModel
}

// Handle can be used to test authentication
// 401 is returned if authentication failed, 200 otherwise
func (a AuthModel) Handle() {
	r := a.GetRequest()
	w := a.GetResponse()

	_, err := auth.HeaderAuth(r, a.GetOidcConfig())

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
