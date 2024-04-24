package converter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

type Media struct {
	Reader io.Reader
	Name   string
}

const tmpDir = "tmp"

func ToWebp(input io.Reader) (*Media, error) {
	err := ensureDir(tmpDir)
	if err != nil {
		return nil, err
	}
	format := "webp"
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s.%s", id, format)
	inDir := filepath.Join(".", tmpDir, name)
	cmd := exec.Command("ffmpeg",
		"-f", "image2pipe",
		"-i", "pipe:0",
		"-vf", "scale='min(1920,iw)':'min(1080,ih)':force_original_aspect_ratio=decrease",
		"-compression_level", "6",
		"-quality", "80",
		"-f", format,
		// i'm not sure why, but the pipe:1 produce corrupted files
		inDir,
	)
	var errBuf bytes.Buffer
	cmd.Stdin = input
	cmd.Stderr = &errBuf
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	cmd.Wait()
	reader, err := toBuffer(inDir)
	if err != nil {
		return nil, err
	}
	return &Media{Reader: reader, Name: name}, err
}

func ensureDir(name string) error {
	err := os.MkdirAll(name, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func toBuffer(name string) (io.Reader, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	err = os.Remove(name)
	r := bytes.NewReader(data)
	return r, err
}
