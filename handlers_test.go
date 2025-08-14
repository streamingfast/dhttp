package dhttp

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httpResponse struct {
	statusCode int
	body       string
}

func TestJSONHandler(t *testing.T) {
	type args struct {
		processor JSONHandlerProcessor
	}
	tests := []struct {
		name string
		args args
		want httpResponse
	}{
		{
			"empty body",
			args{
				processor: func(r *http.Request) (out any, err error) {
					return EmptyBody(), nil
				},
			},
			httpResponse{
				statusCode: http.StatusOK,
				body:       "",
			},
		},
		{
			"body",
			args{
				processor: func(r *http.Request) (out any, err error) {
					return map[string]string{"key": "value"}, nil
				},
			},
			httpResponse{
				statusCode: http.StatusOK,
				body:       `{"key":"value"}` + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", bytes.NewBuffer(nil))
			w := httptest.NewRecorder()

			JSONHandler(tt.args.processor).ServeHTTP(w, req)
			body, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want, httpResponse{
				statusCode: w.Result().StatusCode,
				body:       string(body),
			})
		})
	}
}
