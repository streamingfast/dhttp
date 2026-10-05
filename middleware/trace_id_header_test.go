package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAddTraceIDHeaderMiddleware(t *testing.T) {
	handler := NewTracingLoggingMiddleware(zap.NewNop())(
		NewAddTraceIDHeaderMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})),
	)

	t.Run("generates a trace ID when none is propagated", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

		assert.Len(t, recorder.Header().Get("X-Trace-ID"), 32)
	})

	t.Run("reuses the propagated trace ID", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		assert.Equal(t, "4bf92f3577b34da6a3ce929d0e0e4736", recorder.Header().Get("X-Trace-ID"))
	})
}

func TestAddTraceIDHeaderMiddlewareWithoutTrace(t *testing.T) {
	handler := NewAddTraceIDHeaderMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Empty(t, recorder.Header().Get("X-Trace-ID"))
}
