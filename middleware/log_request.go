package middleware

import (
	"github.com/gorilla/mux"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"net/http"
)

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
