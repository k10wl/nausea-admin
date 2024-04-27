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

	ah := handlers.NewAboutHandler(s.db, t, s.storage)
	mux.HandleFunc("/about/", ah.GetAboutPage)
	mux.HandleFunc("PATCH /about/", ah.PatchAbout)

	fh := handlers.NewFoldersHandler(*s.db, t)
	mux.HandleFunc("/folders/", fh.GetFoldersPage)
	mux.HandleFunc("/folders/{id}", fh.GetFoldersPage)
	mux.HandleFunc("POST /folders/{id}", fh.CreateFolder)
	mux.HandleFunc("PATCH /folders/{id}", fh.PatchFolder)
	mux.HandleFunc("PATCH /folders/{id}/hide", fh.MarkFolderAsDeleted)
	mux.HandleFunc("PATCH /folders/{id}/restore", fh.RestoreFolder)
	mux.HandleFunc("PATCH /folders/{id}/{media_id}", fh.EditFolderMedia)
	mux.HandleFunc("PATCH /folders/{id}/{media_id}/hide", fh.MarkMediaAsDeletedInFolder)
	mux.HandleFunc("PATCH /folders/{id}/{media_id}/restore", fh.RestoreMediaInFolder)

	mh := handlers.NewMediaHandler(*s.db, t, s.storage, l)
	mux.HandleFunc("POST /media", mh.UploadMedia)

	return http.ListenAndServe(s.addr, loggerMux)
}

/*
TODO:
- add upload feature
- create custom modal element because new folder and upload are basically the same
- create service folders that will have unused images
- add about editing
- replace delete with hide, and add real delete button
- global search by name
*/
