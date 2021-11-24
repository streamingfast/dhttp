package dhttp

import (
	"context"
	"encoding/hex"

	"github.com/streamingfast/logging"
	"go.opencensus.io/trace"
)

var defaultTraceID = "00000000000000000000000000000000"

func init() {
	logging.TestingOverride()
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
