package httphelper

import (
	"encoding/json"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// WriteResponse marshals the interface to json and throws on error
// Errors are logged and should be handled by returning 500
func WriteResponse(r *http.Request, w http.ResponseWriter, v interface{}) (err error) {
	log := logconfig.GetLogReq(r)

	var data []byte

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		data, err = json.MarshalIndent(v, "", "  ")
		log.Trace(string(data))
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		log.WithError(err).Error("Failed to encode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(data); err != nil {
		log.WithError(err).Error("Failed to write response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// ReadRequest Parses the request as JSON into the models parameter
// Errors are logged and should be handled by returning 500
func ReadRequest(r *http.Request, w http.ResponseWriter, v interface{}) (err error) {
	log := logconfig.GetLogReq(r)

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Failed to read request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Tracef("Payload: %s", string(bytes))

	err = json.Unmarshal(bytes, v)
	if err != nil {
		log.WithError(err).Error("Failed to read request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// WriteErrorResponse writes the passed error code to the response
// is tries also to resolve a error message (not implemented)
// The status code is set to 400 BadRequest
// Errors are logged and should be handled by returning 500
func WriteErrorResponse(r *http.Request, w http.ResponseWriter, errorCode int) error {
	w.WriteHeader(http.StatusBadRequest)

	return WriteResponse(r, w, struct {
		ErrorCode int    `json:"errorCode"`
		Message   string `json:"message"`
	}{
		ErrorCode: errorCode,
		Message:   "",
	})
}
