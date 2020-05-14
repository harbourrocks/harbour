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

// JsonResp adds the Content-Type: application/json to the response
func JsonResp(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		handler(writer, request)
	}
}
