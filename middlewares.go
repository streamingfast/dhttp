package dhttp

import (
	"net/http"
	"strings"

	stackdriverPropagation "contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/streamingfast/dtracing"
	"github.com/streamingfast/logging"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func NewOpenCensusMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return &ochttp.Handler{
			Handler:     next,
			Propagation: &stackdriverPropagation.HTTPFormat{},
		}
	}
}

// Deprecated: Use `NewOpenCensusMiddleware()` instead
var AddOpenCensusMiddleware = NewOpenCensusMiddleware()

// NewAddLoggerToContextMiddleware adds a request specific logger containing the `trace_id` of the request
// extracted from the request propagation mechanism (StackDriver OpenCensus Propagation currently
// hard-coded) into the request's context.
//
// HTTP handlers can then easily extract a request specific logger using:
//
// `logging.Logger(request.Context(), zlog)`
//
func NewAddLoggerToContextMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return dtracing.NewAddTraceIDAwareLoggerMiddleware(next, logger, &stackdriverPropagation.HTTPFormat{})
	}
}

// Deprecated: Use `NewAddLoggerToContextMiddleware()` instead
var AddLoggerMiddleware = NewAddLoggerToContextMiddleware(zlog)

func NewAddTraceIDHeaderMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			span := trace.FromContext(ctx)
			if span == nil {
				logging.Logger(ctx, logger).Panic("trace is not present in request but should have been")
			}

			w.Header().Set("X-Trace-ID", span.SpanContext().TraceID.String())

			next.ServeHTTP(w, r)
		})
	}
}

// Deprecated: Use `NewAddTraceIDHeaderMiddleware()` instead
var AddTraceMiddleware = NewAddTraceIDHeaderMiddleware(zlog)

// NewLogRequestMiddleware logs important debugging information about the incoming request
// like its `method`, the `host`, the `url` and the `headers`.
func NewLogRequestMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logging.Logger(r.Context(), logger).Debug("handling HTTP request",
				zap.String("method", r.Method),
				zap.Any("host", r.Host),
				zap.Any("url", r.URL),
				zap.Any("headers", r.Header),
			)

			next.ServeHTTP(w, r)
		})
	}
}

// Deprecated: Use `NewLogRequestMiddleware()` instead
var LogMiddleware = NewLogRequestMiddleware(zlog)

func NewCORSMiddleware(allowedOrigins string) mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins(strings.Split(allowedOrigins, ",")),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	)
}
