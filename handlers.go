package dhttp

import (
	"io"
	"net/http"
)

type emptyBody struct{}

// EmptyBody can be returned by a JSONHandlerProcessor to signal the JSON handler
// that the response body should be empty.
//
// This must be used and not `nil` to avoid the JSON handler to write a `null` value.
func EmptyBody() interface{} {
	return emptyBody{}
}

type JSONHandlerProcessor = func(r *http.Request) (out interface{}, err error)

// JSONHandler wraps a simpler `func(r *http.Request) (out interface{}, err error)`
// processor.
//
// If the processor returns something as the `out` value, the `out`
// is serialized as JSON and return to the user.
//
// If the processor returns an error insteand, the `err`
// value is written to the user using `dhttp.WriteError` call.
//
// To have nothing returns as the body of the response, you can do:
//
//	return derr.EmptyBody(), nil
//
// Which will return a 200 OK with an empty body.
func JSONHandler(processor JSONHandlerProcessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := processor(r)
		if err != nil {
			WriteError(r.Context(), w, err)
			return
		}

		if _, ok := out.(emptyBody); ok {
			w.WriteHeader(http.StatusOK)
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

type DirectHandlerProcessor = func(w http.ResponseWriter, r *http.Request) (err error)

// DirectHandler wraps a simpler `func(w http.ResponseWriter, r *http.Request) (err error)`
// processor.
//
// This handler gives you full control over the response writer but with the added
// benefinit of having the error handling done for you.
func DirectHandler(processor DirectHandlerProcessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := processor(w, r)
		if err != nil {
			WriteError(r.Context(), w, err)
			return
		}
	})
}
