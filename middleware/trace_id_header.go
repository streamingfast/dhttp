package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	sftracing "github.com/streamingfast/sf-tracing"
)

// NewAddTraceIDHeaderMiddleware sets the `X-Trace-ID` response header to the trace ID of the
// request. It must run inside `NewTracingLoggingMiddleware`, which is what puts the trace ID
// in the request context; requests without a trace ID get no header.
func NewAddTraceIDHeaderMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if traceID := sftracing.GetTraceID(r.Context()); traceID.IsValid() {
				w.Header().Set("X-Trace-ID", traceID.String())
			}

			next.ServeHTTP(w, r)
		})
	}
}
