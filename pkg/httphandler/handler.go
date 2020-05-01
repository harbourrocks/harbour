package httphandler

import (
	"encoding/json"
	l "github.com/sirupsen/logrus"
	"net/http"
)

type HttpHandler struct {
	Request  *http.Request
	Response http.ResponseWriter
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
