package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"

	"github.com/oojoseph67/ecommerce/internal/utils/interfaces"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{
		provider: provider,
	}
}

func (s *UploadService) UploadProductImage(productId uint, file *multipart.FileHeader) (string, error) {

	extension := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageExtension(extension) {
		return "", fmt.Errorf("invalid file type: %s", extension)
	}

	uploadPath := fmt.Sprintf("products/%d/%s", productId, file.Filename)

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
