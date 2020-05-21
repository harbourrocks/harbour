package handler

import (
	"github.com/harbourrocks/harbour/pkg/logconfig"
	"io/ioutil"
	"net/http"
)

func LogIncoming(w http.ResponseWriter, r *http.Request) {
	log := logconfig.GetLogReq(r)
	log.Trace(r)
	log.Trace(r.URL.String())
	bytes, _ := ioutil.ReadAll(r.Body)
	s := string(bytes)
	log.Trace(s)
}
