package dhttp

import "net/http"

type JSONHandlerProcessor = func(r *http.Request) (out interface{}, err error)

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
