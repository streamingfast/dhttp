package dhttp

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/streamingfast/derr"
	"github.com/streamingfast/validator"
	"github.com/stretchr/testify/assert"
)

func Test_ExtractRequest(t *testing.T) {
	type info struct {
		Prefix string `schema:"prefix"`
		Count  int    `schema:"count"`
		JSON   bool   `schema:"json"`
	}

	r := httptest.NewRequest("GET", "/?prefix=p&count=1&json=true", nil)
	ctx := newTestContext(r.Context())

	request := &info{}
	err := ExtractRequest(ctx, r, request, NewRequestValidator(validator.Rules{
		"prefix": []string{"required"},
		"count":  []string{"min:4"},
	}))

	assert.Equal(t, derr.RequestValidationError(ctx, url.Values{
		"count": []string{"The count field value can not be less than 4"},
	}), err)

	assert.Equal(t, &info{
		Prefix: "p",
		Count:  1,
		JSON:   true,
	}, request)
}

func Test_ExtractRequest_CustomTag(t *testing.T) {
	type info struct {
		Prefix string `json:"prefix"`
	}

	r := httptest.NewRequest("GET", "/", nil)
	ctx := newTestContext(r.Context())

	request := &info{}
	err := ExtractRequest(ctx, r, request, NewRequestValidator(validator.Rules{
		"prefix": []string{"required"},
	}, validator.TagIdentifierOption("json")))

	assert.Equal(t, derr.RequestValidationError(ctx, url.Values{
		"prefix": []string{"The prefix field is required"},
	}), err)

	assert.Equal(t, &info{
		Prefix: "",
	}, request)
}

func Test_ExtractJSONRequest(t *testing.T) {
	type info struct {
		Prefix string `json:"prefix"`
		Count  int    `json:"count"`
		JSON   bool   `json:"json"`
	}

	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"prefix":"p","count":1,"json":true}`))
	ctx := newTestContext(r.Context())

	request := &info{}
	err := ExtractJSONRequest(ctx, r, request, NewJSONRequestValidator(validator.Rules{
		"prefix": []string{"required"},
		"count":  []string{"min:4"},
	}))

	assert.Equal(t, derr.RequestValidationError(ctx, url.Values{
		"count": []string{"The count field value can not be less than 4"},
	}), err)

	assert.Equal(t, &info{
		Prefix: "p",
		Count:  1,
		JSON:   true,
	}, request)
}
