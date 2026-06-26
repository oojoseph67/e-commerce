package interfaces

import (
	"mime/multipart"
)

type UploadProvider interface {
	UploadFile(file *multipart.FileHeader, path string) (url, altText string, err error)
	DeleteFile(path string) error
}
