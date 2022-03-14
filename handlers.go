package dhttp

import (
	"io"
	"net/http"
)

type JSONHandlerProcessor = func(r *http.Request) (out interface{}, err error)

// JSONHandler wraps a simpler `func(r *http.Request) (out interface{}, err error)`
// processor.
//
// If the processor returns something as the `out` value, the `out`
// is serialized as JSON and return to the user.
//
// If the processor returns an error insteand, the `err`
// value is written to the user using `dhttp.WriteError` call.
func JSONHandler(processor JSONHandlerProcessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := processor(r)
		if err != nil {
			WriteError(r.Context(), w, err)
			return
		}

		WriteJSON(r.Context(), w, out)
	})
}

type RawHandlerProcessor = func(r *http.Request) (out io.ReadCloser, err error)

// RawHandler wraps a simpler `func(r *http.Request) (out io.ReadCloser, err error)`
// processor.
//
// If the processor returns something as the `out` value, the `out`
// reader is fully transmitted to the user then close.
//
// If the processor returns an error insteand, the `err`
// value is written to the user using `dhttp.WriteError` call.
func RawHandler(processor RawHandlerProcessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := processor(r)
		if out != nil {
			// Shall we log error if the body cannot be closed properly?
			defer out.Close()
		}

		if err != nil {
			WriteError(r.Context(), w, err)
			return
		}

		WriteFromReader(r.Context(), w, out)
	})
}
