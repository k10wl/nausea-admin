package handlers

import (
	"io"
	"mime/multipart"
	"net/http"

	"nausea-admin/internal/converter"
	"nausea-admin/internal/models"

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
	media, err := converter.ToWebp(file)
	if err != nil {
		errChan <- err
		return
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
	urlChan <- urlWithMediaSize{
		URL: url,
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
