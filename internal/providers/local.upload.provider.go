package providers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

type LocalUploadProvider struct {
	basePath string
	logger   zerolog.Logger
}

func NewLocalUploadProvider(basePath string, logger zerolog.Logger) *LocalUploadProvider {
	return &LocalUploadProvider{
		basePath: basePath,
		logger:   logger,
	}
}

func (up *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (url, altText string, err error) {

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	path = filepath.Join(filepath.Dir(path), newFileName)
	fullPath := filepath.Join(up.basePath, path)

	// creating path like ./upload/2/path
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		up.logger.Err(err).Msg("couldnt generate upload path")
		return "", "", err
	}

	// open file source
	src, err := file.Open()
	if err != nil {
		up.logger.Err(err).Msg("couldnt read file")
		return "", "", err
	}
	defer src.Close()

	// create destination
	dst, err := os.Create(fullPath)
	if err != nil {
		up.logger.Err(err).Msg("couldnt copy file")
		return "", "", err
	}
	defer dst.Close()

	// read from source to destination
	if _, err := dst.ReadFrom(src); err != nil {
		up.logger.Err(err).Msg("couldnt move file to destination")
		return "", "", err
	}

	message := fmt.Sprintf("/uploads/%s", path)

	return message, file.Filename, nil
}

func (up *LocalUploadProvider) DeleteFile(path string) error {
	fullPath := filepath.Join(up.basePath, path)
	return os.Remove(fullPath)
}
