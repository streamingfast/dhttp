package dhttp

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/streamingfast/derr"
)

func FowardResponse(ctx context.Context, w http.ResponseWriter, response *http.Response) {
	// FIXME: Implement using a Pipe stream insteading of reading the full content in memory
	content, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		WriteError(ctx, w, derr.Wrap(err, "unable to read response body while forwarding response"))
	}

	w.WriteHeader(response.StatusCode)
	_, err = w.Write(content)
	if err != nil {
		logWriteResponseError(ctx, "failed forwarding response", err)
	}
}
