package dhttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
