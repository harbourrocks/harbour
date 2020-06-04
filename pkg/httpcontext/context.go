package httpcontext

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const ReqIdKey = "reqId"
const ReqIdHeaderName = "X-Req-Id"

func GetReqIdReq(r *http.Request) string {
	return GetReqIdCtx(r.Context())
}

func GetReqIdCtx(ctx context.Context) string {
	reqId := ctx.Value(ReqIdKey).(string)
	return reqId
}

func UseRequestId(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// todo: build a cors middleware ...
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:4200")

		var reqId string
		if reqId = r.Header.Get(ReqIdHeaderName); reqId == "" {
			// first occurrence
			reqId = uuid.New().String()
		}

		ctx = context.WithValue(ctx, ReqIdKey, reqId)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}

func UseJsonResponse(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}
