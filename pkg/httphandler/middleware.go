package httphandler

import (
	"errors"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
	"net/http"
)

func ForceAuthenticated(model traits.IdTokenTrait) (err error) {
	if model.GetToken() == nil {
		r := model.GetResponse()
		r.WriteHeader(http.StatusUnauthorized)
		err = errors.New("StatusUnauthorized")
	}

	return
}
