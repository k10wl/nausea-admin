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
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", logger(allowGET(handleAboutPage(s))))
	http.HandleFunc("/about/bio", logger(allowGET(handleAboutBio(s))))
	http.HandleFunc("/about/update", logger(allowPOST(handleAboutUpdate(s))))

	// holy-moly this api endpoints look horrific AF
	http.HandleFunc("/contacts", logger(allowGET(handleContactsPage(s))))
	http.HandleFunc("/contacts/data", logger(allowGET(handleContactsData(s))))
	http.HandleFunc("/contacts/update/email", logger(allowPOST(handleEmailUpdate(s))))
	http.HandleFunc("/contacts/update/link", logger(allowPOST(handleLinkUpdate(s))))
	http.HandleFunc("/contacts/delete/link", logger(allowPOST(handleLinkDelete(s))))

	http.HandleFunc("/gallery", logger(allowGET(handleGalleryPage(s))))
	http.HandleFunc("/gallery/upload", logger(allowPOST(handleGalleryUpload(s))))

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
