package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	context2 "github.com/harbourrocks/harbour/pkg/context"
	l "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Get issues a GET request against the url.
//The response is unmarshalled into response
func Get(hRock context2.HRock, url string, response interface{}, token string) (resp *http.Response, err error) {
	log := hRock.L

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

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		l.
			WithError(err).
			WithField("method", req.Method).
			WithField("url", url).
			Error("Failed to send request")
		return
	}

	err = handleResponse(resp, response)
	return
}

// Post issues a POST request against the url.
// The POST payload is specified by body. If body is nil then no body is sent at all.
// The response is unmarshalled into response.
func Post(url string, response interface{}, _ interface{}) (resp *http.Response, err error) {
	resp, err = http.Post(url, "application/json", nil)
	if err != nil {
		l.
			WithError(err).
			WithField("Method", resp.Request.Method).
			WithField("Url", url).
			Error("Failed to send request")
		return
	}

	err = handleResponse(resp, response)
	return
}

func handleResponse(resp *http.Response, response interface{}) (err error) {
	// setup logger with url and status code
	log := l.
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
		log.WithError(err).Error("Failed to unmarshal response: %s", responseBytes)
		return
	} else {
		log.Tracef("Response: %s", responseBytes)
		return
	}
}
