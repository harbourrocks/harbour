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

// NewHttpHandler creates a new HttpHandler and initializes all fields
// this is a shortcut for manually creating the struct
func NewHttpHandler(r *http.Request, w http.ResponseWriter, o auth.OIDCConfig) HttpHandler {
	return HttpHandler{
		Request:    r,
		Response:   w,
		OIDCConfig: o,
	}
}

// WriteResponse marshals the interface to json and throws on error
// Errors are logged and should be handled by returning 500
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

// ReadRequest Parses the request as JSON into the model parameter
// Errors are logged and should be handled by returning 500
func (h *HttpHandler) ReadRequest(model interface{}) (err error) {
	bytes, err := ioutil.ReadAll(h.Request.Body)
	if err != nil {
		l.WithError(err).Error("Failed to read request")
		http.Error(h.Response, err.Error(), http.StatusInternalServerError)
		return
	}

	l.Tracef("Payload: %s", string(bytes))

	err = json.Unmarshal(bytes, model)
	if err != nil {
		l.WithError(err).Error("Failed to read request")
		http.Error(h.Response, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// WriteErrorResponse writes the passed error code to the response
// is tries also to resolve a error message (not implemented)
// The status code is set to 400 BadRequest
// Errors are logged and should be handled by returning 500
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
// Errors are logged and should be handled by returning 401
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
