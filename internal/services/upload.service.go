package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
)

func (s *UploadService) UploadProductImage(productId string, file *multipart.FileHeader) (url, altText string, err error) {

	extension := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExtension(extension) {
		return "", "", fmt.Errorf("invalid file type: %s", extension)
	}

	uploadPath := fmt.Sprintf("products/%s/%s", productId, file.Filename)

	return s.provider.UploadFile(file, uploadPath)
}

func isValidImageExtension(extension string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	// for _, ext := range validExtensions {
	// 	if extension == ext {
	// 		return true
	// 	}
	// }

	// return false

	return slices.Contains(validExtensions, extension)
}
