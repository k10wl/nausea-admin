package handlers

import (
	"fmt"
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/models"
	"nausea-admin/internal/server/logger"
	"nausea-admin/internal/server/template"
	"nausea-admin/internal/server/utils"
	"nausea-admin/internal/storage"
)

type MediaHandler struct {
	DB       db.DB
	Template template.Template
	Storage  *storage.Storage
	Logger   logger.ServerLogger
}

type urlWithMediaSize struct {
	URL          string
	ThumbnailURL string
	models.MediaSize
}

func NewMediaHandler(
	db db.DB,
	t template.Template,
	s *storage.Storage,
	logger logger.ServerLogger,
) MediaHandler {
	return MediaHandler{DB: db, Template: t, Storage: s, Logger: logger}
}

func (mh MediaHandler) UploadMedia(w http.ResponseWriter, r *http.Request) {
	err := parseMultipartForm(r)
	urlValues := r.URL.Query()
	folderId := urlValues.Get("folder-id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "no files found", http.StatusBadRequest)
		return
	}
	urls, errs := filesIntoBucket(files, mh.Storage.AddObject)
	if len(errs) > 0 {
		mh.Logger.Logln(fmt.Sprintf("failed to put into bucket %+v", errs))
	}
	failed := []string{}
	dbDocs := []models.Media{}
	for _, url := range urls {
		media, err := models.NewMedia(url.URL, url.MediaSize)
		media.ThumbnailURL = url.ThumbnailURL
		if err != nil {
			failed = append(failed, url.URL)
			continue
		}
		err = mh.DB.CreateMedia(media)
		if err != nil {
			failed = append(failed, url.URL)
			continue
		}
		dbDocs = append(dbDocs, media)
	}
	if len(failed) > 0 {
		mh.Logger.Logln(fmt.Sprintf("failed to put into DB %+v", failed))
	}
	asContent := []models.MediaContent{}
	for _, m := range dbDocs {
		c, err := m.AsContent(folderId)
		if err != nil {
			mh.Logger.Logln(
				fmt.Sprintf("failed to convert media with id %q to content", c.ID.ID),
			)
			continue
		}
		asContent = append(asContent, c)
	}
	err = mh.DB.UploadMediaToFolder(asContent, folderId)
	if err != nil {
		w.Header().Set("HX-Reswap", "innerHTML")
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	mh.Template.ExecuteTemplate(w, "media-list-range",
		map[string]interface{}{
			"MediaContents": asContent,
		},
	)
}
