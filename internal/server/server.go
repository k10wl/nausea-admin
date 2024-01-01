package server

import (
	"html/template"
	"log"
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	addr    string
	db      *db.DB
	t       *template.Template
	storage *storage.Storage
}

func NewServer(addr string, db *db.DB, t *template.Template, storage *storage.Storage) *Server {
	return &Server{
		addr:    addr,
		db:      db,
		t:       t,
		storage: storage,
	}
}

func (s *Server) Run() error {
	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.Use(logger)

	r.HandleFunc("/", handleAboutPage(s)).Methods(http.MethodGet)
	r.HandleFunc("/", handleAboutUpdate(s)).Methods(http.MethodPost)
	r.HandleFunc("/lazy", handleAboutBio(s)).Methods(http.MethodGet)

	r.HandleFunc("/contacts", handleContactsPage(s)).Methods(http.MethodGet)
	r.HandleFunc("/contacts/lazy", handleContactsLazy(s)).Methods(http.MethodGet)
	r.HandleFunc("/contacts/email", handleEmailPatch(s)).Methods(http.MethodPatch)
	r.HandleFunc("/contacts/links", handleLinkPost(s)).Methods(http.MethodPost)
	r.HandleFunc("/contacts/links/{id}", handleLinkPut(s)).Methods(http.MethodPut)
	r.HandleFunc("/contacts/links/{id}", handleLinkDelete(s)).Methods(http.MethodDelete)

	r.HandleFunc("/gallery", handleGalleryPage(s)).Methods(http.MethodGet)
	// TODO rework this handler, bad url
	r.HandleFunc("/gallery/upload", handleGalleryUpload(s)).Methods(http.MethodGet)

	http.Handle("/", r)
	log.Printf("Listening and serving on %s\n", s.addr)

	return http.ListenAndServe(s.addr, nil)
}

func (s *Server) executeTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := s.t.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
