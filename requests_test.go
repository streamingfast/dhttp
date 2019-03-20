package dhttp

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eoscanada/derr"
	"github.com/eoscanada/validator"
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
