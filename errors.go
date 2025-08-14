package dhttp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/streamingfast/derr"
	"github.com/streamingfast/logging"
	tracing "github.com/streamingfast/sf-tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var zlog, _ = logging.PackageLogger("dhttp", "github.com/streamingfast/dhttp")

// ToErrorResponse turns a plain `error` interface into a proper `ErrorResponse`
// object. It does so with the following rules:
//
// - If `err` is already an `ErrorResponse`, turns it into such and returns it.
// - If `err` was wrapped, find the most cause which is an `ErrorResponse` and returns it.
// - If `err` is a status.Status (or one that was wrapped), convert it to an ErrorResponse
// - Otherwise, return an `UnexpectedError` with the cause sets to `err` received.
func ToErrorResponse(ctx context.Context, err error) *ErrorResponse {
	response := derr.Find(err, isErrorResponse)
	if response != nil {
		return response.(*ErrorResponse)
	}

	response = derr.Find(err, isStatusCode)
	if response != nil {
		status, _ := status.FromError(err)
		return convertStatusToErrorResponse(ctx, status)
	}

	return UnexpectedError(ctx, err)
}

func isStatusCode(err error) bool {
	if _, ok := status.FromError(err); ok {
		return true
	}

	return false
}

func isErrorResponse(err error) bool {
	if _, ok := err.(*ErrorResponse); ok {
		return true
	}

	return false
}

func convertStatusToErrorResponse(ctx context.Context, st *status.Status) *ErrorResponse {
	switch st.Code() {
	case codes.InvalidArgument:
		return BadRequestError(ctx, nil, "request_validation_error", st.Message())
	case codes.Unavailable:
		return ServiceUnavailableError(ctx, nil, "service_unavailable_error", "Service Unavailable")
	case codes.NotFound:
		return NotFoundError(ctx, nil, "not_found_error", st.Message())
	default:
		return UnexpectedError(ctx, st.Err())
	}
}

// Client Errors

func InvalidJSONError(ctx context.Context, err error) *ErrorResponse {
	return BadRequestError(ctx, err, "invalid_json_error", "The request is not a valid json.", "errors", map[string]interface{}{
		"source": err.Error(),
	})
}

func MissingBodyError(ctx context.Context) *ErrorResponse {
	return BadRequestError(ctx, nil, "missing_body_error", "The request body is missing.")
}

func RequestValidationError(ctx context.Context, errors url.Values) *ErrorResponse {
	return BadRequestError(ctx, nil, "request_validation_error", "The request is invalid.", "errors", errors)
}

// Server Errors

// ServiceNotAvailableError represents a failure at the transport layer to reach a given micro-service.
// Note that while `serviceName` is required, it's not directly available to final response for now,
// will probably encrypt it into an opaque string if you ever make usage of it
func ServiceNotAvailableError(ctx context.Context, cause error, serviceName string) *ErrorResponse {
	return ServiceUnavailableError(ctx, cause, "service_unavailable", "The service your are requesting is not currently available.")
}

func UnexpectedError(ctx context.Context, cause error) *ErrorResponse {
	return InternalServerError(ctx, cause, "unexpected_error", "An unexpected error occurred.")
}

// Generic Request Error Classes (4XX)

var (
	BadRequestError                   = newErrorClass(http.StatusBadRequest)
	UnauthorizedError                 = newErrorClass(http.StatusUnauthorized)
	PaymentRequiredError              = newErrorClass(http.StatusPaymentRequired)
	ForbiddenError                    = newErrorClass(http.StatusForbidden)
	NotFoundError                     = newErrorClass(http.StatusNotFound)
	MethodNotAllowedError             = newErrorClass(http.StatusMethodNotAllowed)
	NotAcceptableError                = newErrorClass(http.StatusNotAcceptable)
	ProxyAuthRequiredError            = newErrorClass(http.StatusProxyAuthRequired)
	RequestTimeoutError               = newErrorClass(http.StatusRequestTimeout)
	ConflictError                     = newErrorClass(http.StatusConflict)
	GoneError                         = newErrorClass(http.StatusGone)
	LengthRequiredError               = newErrorClass(http.StatusLengthRequired)
	PreconditionFailedError           = newErrorClass(http.StatusPreconditionFailed)
	RequestEntityTooLargeError        = newErrorClass(http.StatusRequestEntityTooLarge)
	RequestURITooLongError            = newErrorClass(http.StatusRequestURITooLong)
	UnsupportedMediaTypeError         = newErrorClass(http.StatusUnsupportedMediaType)
	RequestedRangeNotSatisfiableError = newErrorClass(http.StatusRequestedRangeNotSatisfiable)
	ExpectationFailedError            = newErrorClass(http.StatusExpectationFailed)
	TeapotError                       = newErrorClass(http.StatusTeapot)
	UnprocessableEntityError          = newErrorClass(http.StatusUnprocessableEntity)
	LockedError                       = newErrorClass(http.StatusLocked)
	FailedDependencyError             = newErrorClass(http.StatusFailedDependency)
	UpgradeRequiredError              = newErrorClass(http.StatusUpgradeRequired)
	PreconditionRequiredError         = newErrorClass(http.StatusPreconditionRequired)
	TooManyRequestsError              = newErrorClass(http.StatusTooManyRequests)
	RequestHeaderFieldsTooLargeError  = newErrorClass(http.StatusRequestHeaderFieldsTooLarge)
	UnavailableForLegalReasonsError   = newErrorClass(http.StatusUnavailableForLegalReasons)
)

// Generic Server Error Classes (5XX)

var (
	InternalServerError                = newErrorClass(http.StatusInternalServerError)
	NotImplementedError                = newErrorClass(http.StatusNotImplemented)
	BadGatewayError                    = newErrorClass(http.StatusBadGateway)
	ServiceUnavailableError            = newErrorClass(http.StatusServiceUnavailable)
	GatewayTimeoutError                = newErrorClass(http.StatusGatewayTimeout)
	StatusHTTPVersionNotSupportedError = newErrorClass(http.StatusHTTPVersionNotSupported)
	VariantAlsoNegotiatesError         = newErrorClass(http.StatusVariantAlsoNegotiates)
	InsufficientStorageError           = newErrorClass(http.StatusInsufficientStorage)
	LoopDetectedError                  = newErrorClass(http.StatusLoopDetected)
	NotExtendedError                   = newErrorClass(http.StatusNotExtended)
	NetworkAuthenticationRequiredError = newErrorClass(http.StatusNetworkAuthenticationRequired)
)

// ErrorFromStatus can be used to programmatically route the right status to one of the HTTP error class above
func ErrorFromStatus(status int, ctx context.Context, cause error, errorCode string, message any, keyvals ...any) *ErrorResponse {
	errorClass := statusToHTTPErrorClass[status]
	if errorClass == nil {
		logError(ctx, "unable to retrieved error class from status, falling back to internal server error", nil, zap.Int("status", status))
		errorClass = InternalServerError
	}

	return errorClass(ctx, cause, errorCode, message, keyvals...)
}

var statusToHTTPErrorClass = map[int]errorClass{
	http.StatusBadRequest:                   BadRequestError,
	http.StatusUnauthorized:                 UnauthorizedError,
	http.StatusPaymentRequired:              PaymentRequiredError,
	http.StatusForbidden:                    ForbiddenError,
	http.StatusNotFound:                     NotFoundError,
	http.StatusMethodNotAllowed:             MethodNotAllowedError,
	http.StatusNotAcceptable:                NotAcceptableError,
	http.StatusProxyAuthRequired:            ProxyAuthRequiredError,
	http.StatusRequestTimeout:               RequestTimeoutError,
	http.StatusConflict:                     ConflictError,
	http.StatusGone:                         GoneError,
	http.StatusLengthRequired:               LengthRequiredError,
	http.StatusPreconditionFailed:           PreconditionFailedError,
	http.StatusRequestEntityTooLarge:        RequestEntityTooLargeError,
	http.StatusRequestURITooLong:            RequestURITooLongError,
	http.StatusUnsupportedMediaType:         UnsupportedMediaTypeError,
	http.StatusRequestedRangeNotSatisfiable: RequestedRangeNotSatisfiableError,
	http.StatusExpectationFailed:            ExpectationFailedError,
	http.StatusTeapot:                       TeapotError,
	http.StatusUnprocessableEntity:          UnprocessableEntityError,
	http.StatusLocked:                       LockedError,
	http.StatusFailedDependency:             FailedDependencyError,
	http.StatusUpgradeRequired:              UpgradeRequiredError,
	http.StatusPreconditionRequired:         PreconditionRequiredError,
	http.StatusTooManyRequests:              TooManyRequestsError,
	http.StatusRequestHeaderFieldsTooLarge:  RequestHeaderFieldsTooLargeError,
	http.StatusUnavailableForLegalReasons:   UnavailableForLegalReasonsError,

	http.StatusInternalServerError:           InternalServerError,
	http.StatusNotImplemented:                NotImplementedError,
	http.StatusBadGateway:                    BadGatewayError,
	http.StatusServiceUnavailable:            ServiceUnavailableError,
	http.StatusGatewayTimeout:                GatewayTimeoutError,
	http.StatusHTTPVersionNotSupported:       StatusHTTPVersionNotSupportedError,
	http.StatusVariantAlsoNegotiates:         VariantAlsoNegotiatesError,
	http.StatusInsufficientStorage:           InsufficientStorageError,
	http.StatusLoopDetected:                  LoopDetectedError,
	http.StatusNotExtended:                   NotExtendedError,
	http.StatusNetworkAuthenticationRequired: NetworkAuthenticationRequiredError,
}

func logError(ctx context.Context, message string, err error, fields ...zap.Field) {
	logging.Logger(ctx, zlog).Error(message, append(fields, zap.Error(err))...)
}

type ErrorResponse struct {
	Code    string         `json:"code"`
	TraceID string         `json:"trace_id"`
	Status  int            `json:"-"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
	Causer  error          `json:"-"`
}

func (e *ErrorResponse) Unwrap() error {
	if e.Causer == nil {
		return nil
	}
	return e.Causer
}

func (e *ErrorResponse) Cause() error { return e.Causer }

func (e *ErrorResponse) ResponseStatus() int { return e.Status }

func (e *ErrorResponse) Error() string {
	index := 0
	details := make([]string, len(e.Details))

	for k, v := range e.Details {
		details[index] = fmt.Sprintf("%s: %v", k, v)
		index++
	}

	detailsString := ""
	if len(details) > 0 {
		detailsString = fmt.Sprintf(" {%s}", strings.Join(details, ", "))
	}

	causeString := ""
	if e.Causer != nil {
		causeString = fmt.Sprintf(" (%s)", e.Causer)
	}

	return fmt.Sprintf("[%s] %d: %s%s%s", e.Code, e.Status, e.Message, causeString, detailsString)
}

type errorClass func(ctx context.Context, cause error, code string, message interface{}, keyvals ...interface{}) *ErrorResponse

func newErrorClass(status int) errorClass {
	return func(ctx context.Context, cause error, code string, message interface{}, keyvals ...interface{}) *ErrorResponse {
		var msg string
		switch actual := message.(type) {
		case string:
			msg = actual
		case error:
			msg = actual.Error()
		case fmt.Stringer:
			msg = actual.String()
		default:
			msg = fmt.Sprintf("%v", actual)
		}

		var details map[string]interface{}
		l := len(keyvals)
		if l > 0 {
			details = make(map[string]interface{})
		}

		for i := 0; i < l; i += 2 {
			k := keyvals[i]
			var v interface{} = "MISSING"
			if i+1 < l {
				v = keyvals[i+1]
			}

			details[fmt.Sprintf("%v", k)] = v
		}

		return &ErrorResponse{Code: code, TraceID: tracing.GetTraceID(ctx).String(), Status: status, Message: msg, Details: details, Causer: cause}
	}
}
