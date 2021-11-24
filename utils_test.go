package dhttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRealIP(t *testing.T) {
	tests := []struct {
		name     string
		header   http.Header
		modifier func(r *http.Request)
		expected string
	}{
		{
			"no x-forward-for header",
			http.Header{},
			nil,
			"192.0.2.1",
		},

		{
			"x-forward-for header, 1 IP", http.Header{
				"X-Forwarded-For": []string{"1.1.1.1"},
			},
			nil,
			"192.0.2.1",
		},

		{
			"x-forward-for header, 2 IP", http.Header{
				"X-Forwarded-For": []string{"1.1.1.1,2.2.2.2"},
			},
			nil,
			"1.1.1.1",
		},

		{
			"x-forward-for header, 3 IP", http.Header{
				"X-Forwarded-For": []string{"1.1.1.1,2.2.2.2,3.3.3.3"},
			},
			nil,
			"2.2.2.2",
		},

		{
			"x-forward-for header, 5 IP", http.Header{
				"X-Forwarded-For": []string{"1.1.1.1,2.2.2.2,3.3.3.3,4.4.4.4"},
			},
			nil,
			"3.3.3.3",
		},

		{
			"no x-forward-for, no remote address somehow",
			nil,
			func(r *http.Request) {
				r.RemoteAddr = ""
			},
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", strings.NewReader("test"))
			request.Header = test.header
			if test.modifier != nil {
				test.modifier(request)
			}

			assert.Equal(t, test.expected, RealIP(request))
		})
	}
}

func Test_FowardErrorResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", strings.NewReader("test"))
	response := http.Response{}
	response.Body = ioutil.NopCloser(strings.NewReader("body"))
	response.StatusCode = 512

	FowardResponse(request.Context(), recorder, &response)

	assert.Equal(t, 512, recorder.Code)
	assert.Equal(t, "body", recorder.Body.String())
}
