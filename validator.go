package dhttp

import (
	"net/http"
	"net/url"

	"github.com/streamingfast/validator"
)

type Validator interface {
	validate(r *http.Request, data interface{}) url.Values
}

type NoOpValidator struct{}

func (v *NoOpValidator) validate(r *http.Request, data interface{}) url.Values {
	return nil
}

var NoValidation = &NoOpValidator{}

type RequestValidator struct {
	rules   validator.Rules
	options []validator.Option
}

func NewRequestValidator(rules validator.Rules, options ...validator.Option) *RequestValidator {
	return &RequestValidator{
		rules:   rules,
		options: append([]validator.Option{validator.TagIdentifierOption("schema")}, options...),
	}
}

func NewJSONRequestValidator(rules validator.Rules, options ...validator.Option) *RequestValidator {
	return &RequestValidator{
		rules:   rules,
		options: append([]validator.Option{validator.TagIdentifierOption("json")}, options...),
	}
}

func (v *RequestValidator) validate(r *http.Request, data interface{}) url.Values {
	return validator.ValidateStruct(data, v.rules, v.options...)
}
