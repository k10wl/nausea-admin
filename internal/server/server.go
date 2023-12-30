package server

import (
	"html/template"
	"log"
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/models"
	"nausea-admin/internal/storage"
)

type Server struct {
	addr    string
	db      *db.DB
	t       *template.Template
	storage *storage.Storage
}

type PageData struct {
	PageMeta
	Info models.Info
}

type PageMeta struct {
	ActiveRoute string
	Title       string
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
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", logger(allowGET(handleAboutPage(s))))
	http.HandleFunc("/gallery", logger(allowGET(handleGalleryPage(s))))
	http.HandleFunc("/gallery/upload", logger(allowPOST(handleGalleryUpload(s))))

	http.HandleFunc("/update", handleBioUpdate(s))

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
