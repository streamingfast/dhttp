package middleware

import (
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/gorilla/mux"
	"github.com/streamingfast/logging"
	sftracing "github.com/streamingfast/sf-tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"net/http"
)

func init() {
	// https://pkg.go.dev/github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator#section-readme
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			// Putting the CloudTraceOneWayPropagator first means the TraceContext propagator
			// takes precedence if both the traceparent and the XCTC headers exist.
			gcppropagator.CloudTraceOneWayPropagator{}, // X-Cloud-Trace-Context instead of traceparent
			propagation.TraceContext{},
			RandomTraceGetter{}, // add a random traceID if there is none yet
			propagation.Baggage{},
		))
}

// NewTracingLoggingMiddleware adds a request specific logger containing the `trace_id` of the request
// extracted from the request propagation mechanism
// HTTP handlers can then easily extract a request specific logger using:
//
// `logging.Logger(request.Context(), zlog)`
func NewTracingLoggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(&addTraceIDMiddleware{next: next, logger: logger}, "", otelhttp.WithTracerProvider(otel.GetTracerProvider()))
	}

}

type addTraceIDMiddleware struct {
	// Actual root logger to instrument with request information
	logger *zap.Logger
	next   http.Handler
}

func (h *addTraceIDMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.Stringer("trace_id", sftracing.GetTraceID(ctx)))
	ctx = logging.WithLogger(ctx, logger)
	h.next.ServeHTTP(w, r.WithContext(ctx))
}
