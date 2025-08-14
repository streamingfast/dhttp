package dhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	pkgErrors "github.com/pkg/errors"
	"github.com/streamingfast/logging"
	tracing "github.com/streamingfast/sf-tracing"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWriteError(t *testing.T) {
	const testTraceID = "00000000000000000000000000000001"

	testContext := func() context.Context {
		return tracing.WithTraceID(context.Background(), tracing.NewFixedTraceID(testTraceID))
	}

	errInvalidJSON := func(id string) *ErrorResponse { return InvalidJSONError(testContext(), errors.New(id)) }
	errUnexpected := func(cause error) *ErrorResponse { return UnexpectedError(testContext(), cause) }

	tests := []struct {
		name               string
		err                error
		expectedStatusCode int
		expectedBody       string
	}{
		{"plain standard error", errors.New("test"), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"plain error response", errInvalidJSON("plain error response"), 400, `
			{"code":"invalid_json_error","trace_id":"%s","details":{"errors":{"source":"plain error response"}},"message":"The request is not a valid json."}
		`},

		{"wrapped error, no cause", pkgErrors.Wrap(nil, "test"), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"wrapped error, cause standard error", pkgErrors.Wrap(errors.New("wrapped"), "test"), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"wrapped error, cause error response", pkgErrors.Wrap(errInvalidJSON("wrapped response"), "test"), 400, `
			{"code":"invalid_json_error","trace_id":"%s","details":{"errors":{"source":"wrapped response"}},"message":"The request is not a valid json."}
		`},

		{"wrapped error, nested cause error response", pkgErrors.Wrap(pkgErrors.Wrap(errInvalidJSON("wrapped again"), "source"), "nested"), 400, `
			{"code":"invalid_json_error","trace_id":"%s","details":{"errors":{"source":"wrapped again"}},"message":"The request is not a valid json."}
		`},

		{"wrapped error, unexpected with response clause", errUnexpected(errInvalidJSON("json")), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"wrapped error, unexpected with wrapped response clause", errUnexpected(pkgErrors.Wrap(errInvalidJSON("json"), "nested1")), 500, `
			{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
		`},

		{"wrapped error, wrapped with unexpected with wrapped response clause", pkgErrors.Wrap(errUnexpected(pkgErrors.Wrap(errInvalidJSON("json"), "nested1")), "deep"), 500, `
		{"code":"unexpected_error","trace_id":"%s","message":"An unexpected error occurred."}
	`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := testContext()
			recorder := httptest.NewRecorder()
			logger := zap.NewExample()

			WriteErrorResponse(logging.WithLogger(ctx, logger), recorder, "prefix", test.err)

			assert.Equal(t, test.expectedStatusCode, recorder.Code)

			expectedBody := strings.TrimSpace(test.expectedBody)
			if strings.Count(expectedBody, "%s") >= 1 {
				expectedBody = fmt.Sprintf(test.expectedBody, testTraceID)
			}

			assert.JSONEq(t, expectedBody, recorder.Body.String())
		})
	}
}
