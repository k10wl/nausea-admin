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
- replace delete with hide, and add real delete button
- update table UI
  - render table as a grid
  - add some hover effect to visually divide elements
  - remove disabled button, why it is there if it is not clickable?
- add upload feature
- add about editing
- create service folders that will have unused images
- add file logger
- global search by name
*/
