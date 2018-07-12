package header

import (
	"net/http"
)

// HandleCors is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.
func HandleCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
