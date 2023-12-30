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

func allowGET(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			errorResponse(w, r, http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func allowPOST(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errorResponse(w, r, http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func errorResponse(w http.ResponseWriter, r *http.Request, code int) {
	log.Printf("Error in request: %s %s --- %d: %s", r.Method, r.URL.Path, code, http.StatusText(code))
	http.Error(w, http.StatusText(code), code)
}

func logger(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		handler(w, r)
	}
}
