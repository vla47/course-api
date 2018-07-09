package header

import (
	"net/http"
	// gorilaContext "github.com/gorilla/context"
)

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// chain and ends the request.

func HandleCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD")

		next.ServeHTTP(w, r)
	})
}


