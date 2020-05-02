package httphandler

import (
	"encoding/json"
	"github.com/harbourrocks/harbour/pkg/auth"
	l "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type HttpHandler struct {
	Request    *http.Request
	Response   http.ResponseWriter
	OIDCConfig auth.OIDCConfig
}

// WriteResponse marshals the interface to json and throws on error
func (h *HttpHandler) WriteResponse(v interface{}) (err error) {
	var data []byte

	if l.GetLevel() == l.TraceLevel {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		http.Error(h.Response, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = h.Response.Write(data); err != nil {
		http.Error(h.Response, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func (h *HttpHandler) ReadRequest(v interface{}) (err error) {
	bytes, err := ioutil.ReadAll(h.Request.Body)
	if err != nil {
		l.WithError(err).Error("Failed to read request")
		http.Error(h.Response, err.Error(), http.StatusInternalServerError)
		return
	}

	l.Tracef("Payload: %s", string(bytes))

	err = json.Unmarshal(bytes, v)
	if err != nil {
		l.WithError(err).Error("Failed to read request")
		http.Error(h.Response, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func (h *HttpHandler) WriteErrorResponse(errorCode int) error {
	h.Response.WriteHeader(http.StatusBadRequest)
	return h.WriteResponse(struct {
		ErrorCode int    `json:"errorCode"`
		Message   string `json:"message"`
	}{
		ErrorCode: errorCode,
		Message:   "",
	})
}

// ExtractUser gets the sub claim of the id_token.
// If the token validation fails an error is returned.
func (h *HttpHandler) ExtractUser() (idToken auth.IdToken, err error) {
	token, err := auth.HeaderAuth(h.Request, h.OIDCConfig)
	if err != nil {
		return // error is logged in HeaderAuth
	}

	idToken, err = auth.IdTokenFromToken(token)
	if err != nil {
		return // error is logged in IdTokenFromToken
	} else {
		l.
			WithField("sub", idToken.Subject).
			WithField("preferred_username", idToken.PreferredUsername).
			WithField("email", idToken.Email).
			WithField("name", idToken.Name).
			Trace("Authenticated user")
	}

	return
}
