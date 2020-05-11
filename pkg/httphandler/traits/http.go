package traits

import (
	"encoding/json"
	"github.com/harbourrocks/harbour/pkg/auth"
	l "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// RequestTrait returns a pointer to the current request object
type HttpTrait interface {
	GetRequest() *http.Request
	GetResponse() http.ResponseWriter
	GetOidcConfig() auth.OIDCConfig
	SetHttp(*http.Request, http.ResponseWriter, auth.OIDCConfig)
}

// RequestModel holds the request
type HttpModel struct {
	request    *http.Request
	response   http.ResponseWriter
	oidcConfig auth.OIDCConfig
}

func (m HttpModel) GetRequest() *http.Request {
	return m.request
}

func (m HttpModel) GetResponse() http.ResponseWriter {
	return m.response
}

func (m HttpModel) GetOidcConfig() auth.OIDCConfig {
	return m.oidcConfig
}

func (m *HttpModel) SetHttp(r *http.Request, w http.ResponseWriter, o auth.OIDCConfig) {
	m.request = r
	m.response = w
	m.oidcConfig = o
}

// WriteResponse marshals the interface to json and throws on error
// Errors are logged and should be handled by returning 500
func (m *HttpModel) WriteResponse(v interface{}) (err error) {
	w := m.GetResponse()

	var data []byte

	if l.GetLevel() == l.TraceLevel {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// ReadRequest Parses the request as JSON into the models parameter
// Errors are logged and should be handled by returning 500
func (m *HttpModel) ReadRequest(model interface{}) (err error) {
	r := m.GetRequest()
	w := m.GetResponse()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.WithError(err).Error("Failed to read request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	l.Tracef("Payload: %s", string(bytes))

	err = json.Unmarshal(bytes, model)
	if err != nil {
		l.WithError(err).Error("Failed to read request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// WriteErrorResponse writes the passed error code to the response
// is tries also to resolve a error message (not implemented)
// The status code is set to 400 BadRequest
// Errors are logged and should be handled by returning 500
func (m *HttpModel) WriteErrorResponse(errorCode int) error {
	w := m.GetResponse()
	w.WriteHeader(http.StatusBadRequest)

	return m.WriteResponse(struct {
		ErrorCode int    `json:"errorCode"`
		Message   string `json:"message"`
	}{
		ErrorCode: errorCode,
		Message:   "",
	})
}

func AddHttp(trait HttpTrait, r *http.Request, w http.ResponseWriter, o auth.OIDCConfig) {
	trait.SetHttp(r, w, o)
}
