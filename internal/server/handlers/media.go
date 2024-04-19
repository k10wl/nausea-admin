package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os/exec"

	"nausea-admin/internal/converter"
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

func Converter(f *multipart.File, cmd *exec.Cmd) (*bytes.Reader, error) {
	// Ensure the file is read from the beginning if reused
	if seeker, ok := (*f).(io.Seeker); ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			log.Printf("Error seeking file: %v", err)
			return nil, err
		}
	}

	var stdoutBuf bytes.Buffer
	cmd.Stdin = *f
	cmd.Stdout = &stdoutBuf
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return bytes.NewReader(stdoutBuf.Bytes()), nil
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
		media, err := models.NewMedia(url)
		if err != nil {
			failed = append(failed, url)
			continue
		}
		err = mh.DB.CreateMedia(media)
		if err != nil {
			failed = append(failed, url)
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
) ([]string, []error) {
	type progressError struct {
		i   int
		err error
	}
	urls := []string{}
	errs := []error{}
	errChan := make(chan progressError)
	urlChan := make(chan string)
	for i, fileHeader := range files {
		go func(
			i int,
			fileHeader *multipart.FileHeader,
			errChan chan progressError,
			urlChan chan string,
		) {
			file, err := fileHeader.Open()
			if err != nil {
				errChan <- progressError{i, err}
				return
			}
			defer file.Close()
			reader, name, err := converter.ToWebp(file)
			if err != nil {
				errChan <- progressError{i, err}
				return
			}
			url, err := uploader(reader, name)
			if err != nil {
				errChan <- progressError{i, err}
			}
			urlChan <- url
		}(
			i,
			fileHeader,
			errChan,
			urlChan,
		)
	}
	for {
		select {
		case err := <-errChan:
			errs = append(errs, err.err)
		case url := <-urlChan:
			urls = append(urls, url)
		}
		if len(urls)+len(errs) == len(files) {
			break
		}
	}
	return urls, errs
}
