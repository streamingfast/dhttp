package dhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/streamingfast/derr"
	"github.com/streamingfast/dtracing"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func WriteText(ctx context.Context, w http.ResponseWriter, content string) {
	ctx, span := dtracing.StartSpan(ctx, "write text response")
	defer span.End()

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(content)); err != nil {
		logWriteResponseError(ctx, "failed writing text response", err)
	}
}

func WriteTextf(ctx context.Context, w http.ResponseWriter, format string, arguments ...interface{}) {
	ctx, span := dtracing.StartSpan(ctx, "write text formatted response")
	defer span.End()

	w.Header().Set("Content-Type", "text/plain")
	if _, err := fmt.Fprintf(w, format, arguments...); err != nil {
		logWriteResponseError(ctx, "failed writing text response", err)
	}
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	ctx, span := dtracing.StartSpan(ctx, "write JSON response", "type", fmt.Sprintf("%T", v))
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logWriteResponseError(ctx, "failed encoding JSON response", err)
	}
}

func WriteJSONString(ctx context.Context, w http.ResponseWriter, json string) {
	ctx, span := dtracing.StartSpan(ctx, "write JSON string response")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(json)); err != nil {
		logWriteResponseError(ctx, "failed writing text response", err)
	}
}

func WriteHTML(ctx context.Context, w http.ResponseWriter, htmlTpl *template.Template, data interface{}) {
	ctx, span := dtracing.StartSpan(ctx, "write HTML formatted response")
	defer span.End()

	w.Header().Set("Content-Type", "text/html")

	if err := htmlTpl.Execute(w, data); err != nil {
		logWriteResponseError(ctx, "failed writing HTML response", err)
	}
}

func WriteFromBytes(ctx context.Context, w http.ResponseWriter, bytes []byte) {
	ctx, span := dtracing.StartSpan(ctx, "write from bytes response")
	defer span.End()

	if _, err := w.Write(bytes); err != nil {
		logWriteResponseError(ctx, "unable to write to client", err)
	}
}

func WriteFromReader(ctx context.Context, w http.ResponseWriter, reader io.Reader) {
	ctx, span := dtracing.StartSpan(ctx, "write from reader response")
	defer span.End()

	if _, err := io.Copy(w, reader); err != nil {
		logWriteResponseError(ctx, "unable to copy to client", err)
	}
}

func WriteError(ctx context.Context, w http.ResponseWriter, err error) {
	ctx, span := dtracing.StartSpan(ctx, "write error response", "type", fmt.Sprintf("%T", err))
	defer span.End()

	derr.WriteError(ctx, w, "unable to fullfil request", err)
}

func logWriteResponseError(ctx context.Context, message string, err error) {
	level := zapcore.ErrorLevel
	if derr.IsClientSideNetworkError(err) {
		level = zapcore.DebugLevel
	}

	logging.Logger(ctx, zlog).Check(level, message).Write(zap.Error(err))
}
