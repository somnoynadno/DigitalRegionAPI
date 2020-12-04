package middleware

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var LogPath = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ping" {
			IP := r.Header.Get("X-Real-IP") // depends on nginx
			log.Info(fmt.Sprintf("%s: %s %s (%s)", IP, r.Method, r.RequestURI, r.Host))
		}
		next.ServeHTTP(w, r)
	})
}

var LogBody = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Error reading body: %v", err)
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}

		if len(body) > 0 {
			log.Debug(string(body))
		}

		// And now set a new body, which will simulate the same data we read:
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// Call next handler, passing the response wrapper:
		next.ServeHTTP(w, r)
	})
}
