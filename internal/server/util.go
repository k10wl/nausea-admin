package server

import (
	"log"
	"net/http"
	"time"
)

func errorResponse(w http.ResponseWriter, r *http.Request, code int, e error) {
	log.Printf("Error in request: %s %s --- %d: %s --- Error: %s", r.Method, r.URL.Path, code, http.StatusText(code), e)
	msg := e.Error()
	if msg == "" {
		msg = http.StatusText(code)
	}
	log.Println(msg)
	http.Error(w, msg, code)
}

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Method: %s, URL: %s, Duration: %s\n", r.Method, r.URL.Path, duration)
	})
}
