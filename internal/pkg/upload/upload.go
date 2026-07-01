package upload

import (
	"path/filepath"
	"strings"
	"time"
	"github.com/google/uuid"
)

func GenerateFilename(original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	now := time.Now()
	return filepath.Join(now.Format("2006"), now.Format("01"), uuid.New().String()+ext)
}

func IsImage(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}
