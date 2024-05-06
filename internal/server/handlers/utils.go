package handlers

import (
	"bytes"
	"fmt"
	"io"

	"mime/multipart"
	"net/http"

	"nausea-admin/internal/converter"
	"nausea-admin/internal/models"

	"github.com/google/uuid"
	"golang.org/x/image/webp"
)

func parseMultipartForm(r *http.Request) error {
	return r.ParseMultipartForm(1 << 30) // memory limit of 1GB
}

func processFile(
	fileHeader *multipart.FileHeader,
	uploader func(io.Reader, string) (string, error),
	errChan chan error,
	urlChan chan urlWithMediaSize,
) {
	file, err := fileHeader.Open()
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()
	var buf bytes.Buffer
	tee := io.TeeReader(file, &buf)
	id, err := uuid.NewUUID()
	if err != nil {
		errChan <- err
		return
	}
	uniqName := id.String()
	media, err := converter.ToWebp(
		converter.Opts{
			Input:     tee,
			Name:      uniqName,
			MinWidth:  1920,
			MinHeight: 1080,
			Quality:   80,
		})
	fmt.Printf("media: %v\n", media)
	if err != nil {
		errChan <- err
		return
	}
	thumbnail, err := converter.ToWebp(
		converter.Opts{
			Input:     &buf,
			Name:      uniqName + "-thumbnail",
			MinWidth:  480,
			MinHeight: 360,
			Quality:   60,
		})
	fmt.Printf("thumbnail: %v\n", thumbnail)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		errChan <- err
	}
	img, err := webp.DecodeConfig(media.Reader)
	if err != nil {
		errChan <- err
		return
	}
	url, err := uploader(media.Reader, media.Name)
	if err != nil {
		errChan <- err
	}
	thumbnailURL, err := uploader(thumbnail.Reader, thumbnail.Name)
	if err != nil {
		errChan <- err
	}
	fmt.Printf("thumbnailURL: %v\n", thumbnailURL)
	urlChan <- urlWithMediaSize{
		URL:          url,
		ThumbnailURL: thumbnailURL,
		MediaSize: models.MediaSize{
			Width:  img.Width,
			Height: img.Height,
		},
	}
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

func getFolderID(r *http.Request) string {
	folderID := r.PathValue("id")
	if folderID == "" {
		folderID = models.RootFolderID
	}
	return folderID
}
