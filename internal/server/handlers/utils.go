package handlers

import (
	"io"
	"mime/multipart"

	"nausea-admin/internal/converter"
	"nausea-admin/internal/models"

	"golang.org/x/image/webp"
)

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
