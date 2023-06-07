package middleware

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"strings"
)

func NewCORSMiddleware(allowedOrigins string) mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins(strings.Split(allowedOrigins, ",")),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	)
}
