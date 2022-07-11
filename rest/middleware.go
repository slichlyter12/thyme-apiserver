package rest

import (
	"log"
	"net/http"
	"time"
)

func alwaysJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		log.Printf("%s [%s](%s) @ %s", time.Since(now), r.RequestURI, r.Method, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
