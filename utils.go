package dhttp

import (
	"context"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/streamingfast/derr"
)

var portSuffixRegex = regexp.MustCompile(`:[0-9]{2,5}$`)

// RealIP tries to determine the actual client remote IP that initiated the connection.
// This method is aware of Google Cloud Load Balancer and if the `X-Forwarded-For` exists
// and has 2 or more addresses, it will assume it's coming from Google Cloud Load Balancer.
//
// When behind a Google Load Balancer, the only two values that we can
// be sure about are the `n - 2` and `n - 1` (so the last two values
// in the array). The very last value (`n - 1`) is the Google IP and the
// `n - 2` value is the actual remote IP that reached the load balancer.
//
// When there is more than 2 IPs, all other values prior `n - 2` are
// those coming from the `X-Forwarded-For` HTTP header received by the load
// balancer directly, so something a client might have added manually. Since
// they are coming from an HTTP header and not from Google directly, they
// can be forged and cannot be trusted.
//
// @see https://cloud.google.com/load-balancing/docs/https#x-forwarded-for_header
func RealIP(r *http.Request) string {
	xForwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if xForwardedFor != "" {
		addresses := strings.Split(xForwardedFor, ",")
		if len(addresses) >= 2 {
			return addresses[len(addresses)-2]
		}
	}

	if r.RemoteAddr != "" {
		// The RemoteAddr actually has a format of the form `<ip>:<port>`, we remove the port suffix part
		return portSuffixRegex.ReplaceAllString(r.RemoteAddr, "")
	}

	return "0.0.0.0"
}

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
