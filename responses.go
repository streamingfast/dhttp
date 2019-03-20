package dhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/eoscanada/derr"
	"github.com/eoscanada/dtracing"
	"github.com/eoscanada/logging"
)

func WriteError(ctx context.Context, w http.ResponseWriter, err error) {
	ctx, span := dtracing.StartSpan(ctx, "write error response", "type", fmt.Sprintf("%T", err))
	defer span.End()

	derr.WriteError(ctx, w, "unable to fullfil request", err)
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	ctx, span := dtracing.StartSpan(ctx, "write JSON response", "type", fmt.Sprintf("%T", v))
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logWriteResponseError(ctx, "failed encoding JSON response", err)
	}
}

func logWriteResponseError(ctx context.Context, message string, err error) {
	level := zapcore.ErrorLevel
	if derr.IsClientSideNetworkError(err) {
		level = zapcore.DebugLevel
	}

	logging.Logger(ctx, zlog).Check(level, message).Write(zap.Error(err))
}
