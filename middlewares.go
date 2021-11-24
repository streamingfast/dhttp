package dhttp

import (
	"net/http"
	"strings"

	stackdriverPropagation "contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/streamingfast/logging"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func AddOpenCensusMiddleware(next http.Handler) http.Handler {
	return &ochttp.Handler{
		Handler:     next,
		Propagation: &stackdriverPropagation.HTTPFormat{},
	}
}

func AddLoggerMiddleware(next http.Handler) http.Handler {
	return &logging.Handler{
		Next:        next,
		Propagation: &stackdriverPropagation.HTTPFormat{},
		RootLogger:  zlog,
	}
}

func AddTraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.FromContext(ctx)
		if span == nil {
			logging.Logger(ctx, zlog).Panic("trace is not present in request but should have been")
		}

		w.Header().Set("X-Trace-ID", span.SpanContext().TraceID.String())

		next.ServeHTTP(w, r)
	})
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Logger(r.Context(), zlog).Debug("handling HTTP request",
			zap.String("method", r.Method),
			zap.Any("host", r.Host),
			zap.Any("url", r.URL),
			zap.Any("headers", r.Header),
		)

		next.ServeHTTP(w, r)
	})
}

func NewCORSMiddleware(allowedOrigins string) mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins(strings.Split(allowedOrigins, ",")),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	)
}
