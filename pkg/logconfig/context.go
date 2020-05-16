package logconfig

import (
	"context"
	"github.com/harbourrocks/harbour/pkg/httpcontext"
	"github.com/sirupsen/logrus"
	"net/http"
)

const LoggerKey = "logger"

func GetLogReq(r *http.Request) *logrus.Entry {
	return GetLogCtx(r.Context())
}

func GetLogCtx(ctx context.Context) *logrus.Entry {
	log := ctx.Value(LoggerKey)

	if log == nil {
		logrus.Fatal("Logger is missing in the context") // panics
	}

	return log.(*logrus.Entry)
}

func UseLogger(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if reqId := ctx.Value(httpcontext.ReqIdKey); reqId == nil {
			logrus.Fatal("No request id associated with request") //  panics
		} else {
			log := logrus.WithField(httpcontext.ReqIdKey, reqId)
			ctx = context.WithValue(ctx, LoggerKey, log)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return fn
}
