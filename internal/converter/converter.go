package converter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Media struct {
	Reader io.Reader
	Name   string
}

const tmpDir = "tmp"

type Opts struct {
	Input     io.Reader
	MinWidth  int
	MinHeight int
	Quality   int
	Name      string
}

func ToWebp(opts Opts) (*Media, error) {
	err := ensureDir(tmpDir)
	if err != nil {
		return nil, err
	}
	format := "webp"
	fileName := fmt.Sprintf("%s.%s", opts.Name, format)
	inDir := filepath.Join(".", tmpDir, fileName)
	fmt.Printf("inDir: %v\n", inDir)
	cmd := exec.Command("ffmpeg",
		"-f", "image2pipe",
		"-i", "pipe:0",
		"-vf", fmt.Sprintf(
			"scale='min(%d,iw)':'min(%d,ih)':force_original_aspect_ratio=decrease",
			opts.MinWidth,
			opts.MinHeight,
		),
		"-compression_level", "6",
		"-quality", fmt.Sprint(opts.Quality),
		"-f", format,
		// I'm not sure why, but the pipe:1 produce corrupted files
		inDir,
	)
	var errBuf bytes.Buffer
	cmd.Stdin = opts.Input
	cmd.Stderr = &errBuf
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}
	reader, err := toBuffer(inDir)
	if err != nil {
		return nil, err
	}
	return &Media{Reader: reader, Name: fileName}, err
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
	err = os.Remove(name)
	r := bytes.NewReader(data)
	return r, err
}
