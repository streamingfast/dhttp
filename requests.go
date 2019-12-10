package dhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/eoscanada/derr"
	"github.com/eoscanada/logging"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.RegisterConverter(time.Duration(0), stringToTimeDuration)
}

func ExtractRequest(ctx context.Context, r *http.Request, request interface{}, validator Validator) error {
	err := decoder.Decode(request, requestToSchemaDecodingMap(r))
	if err != nil {
		return sanitizeSchemaError(ctx, err)
	}

	requestErrors := validator.validate(r, request)
	if len(requestErrors) > 0 {
		return derr.RequestValidationError(ctx, requestErrors)
	}

	return nil
}

func ExtractJSONRequest(ctx context.Context, r *http.Request, request interface{}, validator Validator) error {
	if r.Body == nil {
		return derr.MissingBodyError(ctx)
	}

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return derr.InvalidJSONError(ctx, err)
	}

	requestErrors := validator.validate(r, request)
	if len(requestErrors) > 0 {
		return derr.RequestValidationError(ctx, requestErrors)
	}

	return nil
}

func requestToSchemaDecodingMap(r *http.Request) url.Values {
	variables := r.URL.Query()

	pathVariables := mux.Vars(r)
	for key, pathVariable := range pathVariables {
		variables[key] = append(variables[key], pathVariable)
	}

	return variables
}

func sanitizeSchemaError(ctx context.Context, err error) error {
	zlogger := logging.Logger(ctx, zlog)
	errors := url.Values{}

	switch v := err.(type) {
	case schema.MultiError:
		for field, childErr := range v {
			errors[field] = []string{schemErrorToString(ctx, zlogger, childErr)}
		}

	default:
		errors["_global"] = []string{err.Error()}
	}

	return derr.RequestValidationError(ctx, errors)
}

func schemErrorToString(ctx context.Context, zlogger *zap.Logger, err error) string {
	if v, ok := err.(schema.ConversionError); ok {
		if v.Err != nil {
			zlogger.Debug("Conversion underlying error", zap.Error(v.Err))
		}

		return fmt.Sprintf("Unable to convert value to expected type %v", v.Type)
	}

	zlogger.Info("Unknown conversion error", zap.Error(err))
	return "Unknow conversion error, invalid value"
}

func stringToTimeDuration(input string) reflect.Value {
	if duration, err := time.ParseDuration(input); err == nil {
		return reflect.ValueOf(duration)
	}

	return reflect.Value{}
}
