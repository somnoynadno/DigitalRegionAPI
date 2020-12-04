package middleware

import (
	"net/http"
	"strings"
)

var CORS = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			w.Header().Add("Content-Type", "application/json")
		}

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
		}
	})
}