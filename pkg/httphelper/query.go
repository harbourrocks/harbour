package httphelper

import "net/http"

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
