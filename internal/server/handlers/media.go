package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"nausea-admin/internal/db"
	"nausea-admin/internal/models"
	"nausea-admin/internal/server/logger"
	"nausea-admin/internal/server/template"
	"nausea-admin/internal/storage"
)

type MediaHandler struct {
	DB       db.DB
	Template template.Template
	Storage  *storage.Storage
	Logger   logger.ServerLogger
}

type urlWithMediaSize struct {
	URL string
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
	err := r.ParseMultipartForm(1 << 30) // memory limit of 1GB
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
		mh.Logger.Logln(fmt.Sprintf("failed to put into DB %+v", err))
	}
	mh.Template.ExecuteTemplate(w, "media-list-range",
		map[string]interface{}{
			"MediaContents": asContent,
		},
	)
}

func filesIntoBucket(
	files []*multipart.FileHeader,
	uploader func(io.Reader, string) (string, error),
) ([]urlWithMediaSize, []error) {
	urls := []urlWithMediaSize{}
	errs := []error{}
	errChan := make(chan error)
	urlChan := make(chan urlWithMediaSize)
	for _, fileHeader := range files {
		go processFile(
			fileHeader,
			uploader,
			errChan,
			urlChan,
		)
	}
	for {
		select {
		case err := <-errChan:
			errs = append(errs, err)
		case url := <-urlChan:
			urls = append(urls, url)
		}
		if len(urls)+len(errs) == len(files) {
			break
		}
	}
	return urls, errs
}
