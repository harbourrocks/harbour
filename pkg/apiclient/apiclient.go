package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/harbourrocks/harbour/pkg/httpcontext"
	"github.com/harbourrocks/harbour/pkg/logconfig"
	l "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Get issues a GET request against the url.
//The response is unmarshalled into response
func Get(ctx context.Context, url string, response interface{}, token string, header map[string]string) (resp *http.Response, err error) {
	log := logconfig.GetLogCtx(ctx)
	reqId := httpcontext.GetReqIdCtx(ctx)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		log.
			WithField("url", url).WithError(err).
			Error("Failed to create request")
		return
	}

	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	if reqId != "" {
		req.Header.Add(httpcontext.ReqIdHeaderName, reqId)
	}

	if header != nil {
		for name, value := range header {
			req.Header.Add(name, value)
		}
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		l.
			WithError(err).
			WithField("method", req.Method).
			WithField("url", url).
			Error("Failed to send request")
		return
	}

	err = handleResponse(ctx, resp, response)
	return
}

// Post issues a POST request against the url.
// The POST payload is specified by body. If body is nil then no body is sent at all.
// The response is unmarshalled into response.
func Post(ctx context.Context, url string, response interface{}, body interface{}, token string, header map[string]string) (resp *http.Response, err error) {
	log := logconfig.GetLogCtx(ctx)
	reqId := httpcontext.GetReqIdCtx(ctx)

	var req *http.Request
	if body != nil {
		bodyBytes := new(bytes.Buffer)
		err = json.NewEncoder(bodyBytes).Encode(body)
		if err != nil {
			log.WithError(err).Error("serialization of body failed")
			return
		}

		req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, url, bodyBytes)
	} else {
		req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, url, nil)
	}

	if err != nil {
		log.
			WithField("url", url).WithError(err).
			Error("Failed to create request")
		return
	}

	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	if reqId != "" {
		req.Header.Add(httpcontext.ReqIdHeaderName, reqId)
	}

	req.Header.Add("Content-Type", "application/json")

	if header != nil {
		for name, value := range header {
			req.Header.Add(name, value)
		}
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.
			WithError(err).
			WithField("Method", req.Method).
			WithField("Url", url).
			Error("Failed to send request")
		return
	}

	if response != nil {
		err = handleResponse(ctx, resp, response)
	}

	return
}

func handleResponse(ctx context.Context, resp *http.Response, response interface{}) (err error) {
	log := logconfig.GetLogCtx(ctx)

	// setup logger with url and status code
	log = log.
		WithField("Method", resp.Request.Method).
		WithField("Url", resp.Request.URL.String()).
		WithField("StatusCode", resp.StatusCode)

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Failed to receive response")
		return
	}

	if resp.StatusCode >= 400 {
		log.WithError(err).Errorf("Failed status code: %s", responseBytes)
		return
	}

	if err = json.Unmarshal(responseBytes, response); err != nil {
		log.WithError(err).Errorf("Failed to unmarshal response: %s", string(responseBytes))
		return
	} else {
		log.Tracef("Response: %s", string(responseBytes))
		return
	}
}
