package dhttp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zlog, _ = logging.PackageLogger("dhttp", "github.com/streamingfast/dhttp")

// NewLoggingRoundTripper create a wrapping `http.RoundTripper` aware object that intercepts
// the request as well as the response and logs them to the specified logger according to
// some rules if the debug and tracing level are enabled or not.
//
// If the received `next` argument is set as `nil`, the `http.DefaultTransport` value
// will be used as the actual transport handler.
//
// If debug is enabled, a one-line Request log containing HTTP method, URL and Headers is logged,
// and if tracing is enabled, the full request is dumped to the logger, in a multi-line
// log.
//
// If debug is enabled, a one-line Response log containing HTTP status and body length is logged,
// and if tracing is enabled, the full request is dumped to the logger, in a multi-line
// log.
//

func NewLoggingRoundTripper(logger *zap.Logger, tracer logging.Tracer, next http.RoundTripper) *LoggingRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}

	return &LoggingRoundTripper{
		transport: next,
		logger:    logger,
		tracer:    tracer,
	}
}

type LoggingRoundTripper struct {
	transport http.RoundTripper
	logger    *zap.Logger
	tracer    logging.Tracer
}

func (t *LoggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	logger := logging.Logger(request.Context(), t.logger)
	debugEnabled := logger.Core().Enabled(zap.DebugLevel)

	if debugEnabled {
		traceEnabled := t.tracer.Enabled()

		if traceEnabled {
			requestDump, err := httputil.DumpRequestOut(request, true)
			if err != nil {
				logger.Debug(fmt.Sprintf("HTTP request %s %s (unable to log request body: %s)", request.Method, request.URL.String(), err), zap.Array("headers", zapHeaders(request.Header)))
			} else {
				logger.Debug("HTTP request\n" + string(requestDump))
			}
		} else {
			logger.Debug(fmt.Sprintf("HTTP request %s %s", request.Method, request.URL.String()), zap.Array("headers", zapHeaders(request.Header)))
		}
	}

	response, err := t.transport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if debugEnabled {
		traceEnabled := t.tracer.Enabled()

		if traceEnabled {
			responseDump, err := httputil.DumpResponse(response, true)
			if err != nil {
				logger.Debug(fmt.Sprintf("HTTP response %s (%d bytes, unable to log response body: %s)", response.Status, response.ContentLength, err))
			} else {
				logger.Debug("HTTP response\n" + string(responseDump))
			}
		} else {
			logger.Debug(fmt.Sprintf("HTTP response %s (%d bytes)", response.Status, response.ContentLength))
		}
	}

	return response, nil
}

type zapHeaders http.Header

func (h zapHeaders) MarshalLogArray(encoder zapcore.ArrayEncoder) error {
	for key, values := range h {
		encoder.AppendString(fmt.Sprintf("%s: %s", toPascalCase(key), strings.Join(values, " <> ")))
	}

	return nil
}

var upperRegex = regexp.MustCompile("[A-Z]")

func toPascalCase(in string) string {
	raw := strcase.ToCamel(in)

	return strings.TrimSpace(upperRegex.ReplaceAllStringFunc(raw, func(s string) string {
		return " " + s
	}))
}
