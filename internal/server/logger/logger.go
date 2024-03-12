package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type ServerLogger struct{}

func NewServerLogger() ServerLogger {
	return ServerLogger{}
}

func (sl ServerLogger) HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		duration := time.Since(start)
		sl.Log(fmt.Sprintf("Method: %s, URL: %s, Duration: %s\n", r.Method, r.URL.Path, duration))
	})
}

func (sl ServerLogger) Log(s string) {
	log.Print(s)
}
