package server

import (
	"log"
	"net/http"
)

func pageMeta(r *http.Request) PageMeta {
	return PageMeta{
		ActiveRoute: r.RequestURI,
	}
}

func errorResponse(w http.ResponseWriter, r *http.Request, code int, e error) {
	log.Printf("Error in request: %s %s --- %d: %s --- Error: %s", r.Method, r.URL.Path, code, http.StatusText(code), e)
	http.Error(w, http.StatusText(code), code)
}

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}
