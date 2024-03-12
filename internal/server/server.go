package server

import (
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/server/handlers"
	"nausea-admin/internal/server/logger"
	"nausea-admin/internal/server/template"
	"nausea-admin/internal/storage"
)

type Server struct {
	addr    string
	db      *db.DB
	storage *storage.Storage
}

func NewServer(addr string, db *db.DB, storage *storage.Storage) *Server {
	return &Server{
		addr:    addr,
		db:      db,
		storage: storage,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	l := logger.NewServerLogger()
	defer l.CloseLogger()
	loggerMux := l.HTTPLogger(mux)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	t := template.NewTemplate(l)

	hh := handlers.NewHomeHandler(t)
	mux.HandleFunc("/", hh.GetHomePage)

	fh := handlers.NewFoldersHandler(*s.db, t)
	mux.HandleFunc("/folders/", fh.GetFoldersPage)
	mux.HandleFunc("/folders/{id}", fh.GetFoldersPage)
	mux.HandleFunc("POST /folders/{id}", fh.CreateFolder)
	mux.HandleFunc("DELETE /folders/{id}", fh.DeleteFolder)
	mux.HandleFunc("PATCH /folders/{id}/restore", fh.RestoreFolder)
	return http.ListenAndServe(s.addr, loggerMux)
}

/*
TODO:
- add upload feature
- create service folders that will have unused images
- add about editing
- replace delete with hide, and add real delete button
- global search by name
*/
