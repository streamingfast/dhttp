package dhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/streamingfast/derr"
	"github.com/streamingfast/logging"
	tracing "github.com/streamingfast/sf-tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func WriteText(ctx context.Context, w http.ResponseWriter, content string) {
	ctx, span := tracing.GetTracer().Start(ctx, "write text response")
	defer span.End()

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(content)); err != nil {
		logWriteResponseErrorCtx(ctx, "failed writing text response", err)
	}
}

func WriteTextf(ctx context.Context, w http.ResponseWriter, format string, arguments ...any) {
	ctx, span := tracing.GetTracer().Start(ctx, "write text formatted response")
	defer span.End()

	w.Header().Set("Content-Type", "text/plain")
	if _, err := fmt.Fprintf(w, format, arguments...); err != nil {
		logWriteResponseErrorCtx(ctx, "failed writing text response", err)
	}
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, v any) {
	ctx, span := tracing.GetTracer().Start(ctx, "write JSON response",
		trace.WithAttributes(attribute.String("type", fmt.Sprintf("%T", v))),
	)
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logWriteResponseErrorCtx(ctx, "failed encoding JSON response", err)
	}
}

func WriteJSONString(ctx context.Context, w http.ResponseWriter, json string) {
	ctx, span := tracing.GetTracer().Start(ctx, "write JSON string response")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(json)); err != nil {
		logWriteResponseErrorCtx(ctx, "failed writing text response", err)
	}
}

func WriteHTML(ctx context.Context, w http.ResponseWriter, htmlTpl *template.Template, data any) {
	ctx, span := tracing.GetTracer().Start(ctx, "write HTML formatted response")
	defer span.End()

	w.Header().Set("Content-Type", "text/html")

	if err := htmlTpl.Execute(w, data); err != nil {
		logWriteResponseErrorCtx(ctx, "failed writing HTML response", err)
	}
}

func WriteFromBytes(ctx context.Context, w http.ResponseWriter, bytes []byte) {
	ctx, span := tracing.GetTracer().Start(ctx, "write from bytes response")
	defer span.End()

	if _, err := w.Write(bytes); err != nil {
		logWriteResponseErrorCtx(ctx, "unable to write to client", err)
	}
}

func WriteFromReader(ctx context.Context, w http.ResponseWriter, reader io.Reader) {
	ctx, span := tracing.GetTracer().Start(ctx, "write from reader response")
	defer span.End()

	if _, err := io.Copy(w, reader); err != nil {
		logWriteResponseErrorCtx(ctx, "unable to copy to client", err)
	}
}

// Write a generic error response to HTTP, the error is going to transformed to
// a [ErrorResponse] via a call to [ToErrorResponse].
//
// A generic `unable to fullfil request` message will be logged, use [WriteErrorResponse]
// to control which message will be logged.
func WriteError(ctx context.Context, w http.ResponseWriter, err error) {
	ctx, span := tracing.GetTracer().Start(ctx, "write error response", trace.WithAttributes(
		attribute.String("type", fmt.Sprintf("%T", err)),
	))
	defer span.End()

	WriteErrorResponse(ctx, w, "unable to fullfil request", err)
}

// WriteError writes the receiver error to HTTP and log it into a Zap logger at the same
// time with the right level based on the actual status code. The `WriteError` handles
// various type for the `err` parameter.
//
// Deprecated: HTTP error handling has moved to [dhttp](https://github.com/streamingfast/dhttp) package,
// which provides a more comprehensive and flexible approach to error handling in HTTP contexts.
func WriteErrorResponse(ctx context.Context, w http.ResponseWriter, message string, err error) {
	response := ToErrorResponse(ctx, err)
	zlogger := logging.Logger(ctx, zlog)

	if ctx.Err() != context.Canceled && response.ResponseStatus() >= 500 {
		zlogger.Error(message, zap.Error(err))
	} else {
		zlogger.Debug(message, zap.Error(err))
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(response.ResponseStatus())

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logWriteResponseErrorLogger(zlogger, "unable to serialize error response", err)
	}
}

func logWriteResponseErrorCtx(ctx context.Context, message string, err error) {
	logWriteResponseErrorLogger(logging.Logger(ctx, zlog), message, err)
}

func logWriteResponseErrorLogger(logger *zap.Logger, message string, err error) {
	level := zapcore.ErrorLevel
	if derr.IsClientSideNetworkError(err) {
		level = zapcore.DebugLevel
	}

	logger.Check(level, message).Write(zap.Error(err))
}
