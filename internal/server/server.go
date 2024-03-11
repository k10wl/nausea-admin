package server

import (
	"html/template"
	"log"
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/storage"
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
	mux := http.NewServeMux()
	loggerMux := logger(mux)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/", GetHomePage(s))
	mux.HandleFunc("/folders/", GetFoldersPage(s))
	mux.HandleFunc("POST /folders/{id}", CreateFolder(s))
	mux.HandleFunc("DELETE /folders/{id}", DeleteFolder(s))
	mux.HandleFunc("PATCH /folders/{id}/restore", RestoreFolder(s))
	mux.HandleFunc("/folders/{id}", GetFoldersPage(s))
	return http.ListenAndServe(s.addr, loggerMux)
}

func (s *Server) executeTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := s.t.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
