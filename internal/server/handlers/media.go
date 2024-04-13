package handlers

import (
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/server/template"
	"nausea-admin/internal/storage"
)

type MediaHandler struct {
	DB       db.DB
	Template template.Template
	Storage  *storage.Storage
}

func NewMediaHandler(db db.DB, t template.Template, s *storage.Storage) MediaHandler {
	return MediaHandler{DB: db, Template: t, Storage: s}
}

func (mh MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
}
