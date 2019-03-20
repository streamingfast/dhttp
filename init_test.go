package dhttp

import (
	"context"
	"encoding/hex"
	"os"

	"go.opencensus.io/trace"

	"go.uber.org/zap"
)

var defaultTraceID = "00000000000000000000000000000000"

func init() {
	if os.Getenv("DEBUG") != "" {
		zlog, _ = zap.NewDevelopment()
	}
}

func newTestContext(parentCtx context.Context) context.Context {
	traceID := fixedTraceID(defaultTraceID)
	spanContext := trace.SpanContext{TraceID: traceID}
	ctx, _ := trace.StartSpanWithRemoteParent(parentCtx, "test", spanContext)

	return ctx
}

func fixedTraceID(hexInput string) (out trace.TraceID) {
	rawTraceID, _ := hex.DecodeString(hexInput)
	copy(out[:], rawTraceID)

	return
}
