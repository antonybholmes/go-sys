package sys

import (
	"mime"
	"path/filepath"
	"strings"
)

var (
	CommonMimeTypes = map[string]string{
		".json": "application/json",
		".txt":  "text/plain; charset=utf-8",
		".csv":  "text/csv; charset=utf-8",
		".md":   "text/markdown; charset=utf-8",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
	}
)

func DetectContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	if mimeType, ok := CommonMimeTypes[ext]; ok {
		return mimeType
	}

	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}

	return "application/octet-stream"
}

// DetectFileExtension returns the file extension for a given filename.
// If the extension is a common MIME type, it returns the extension.
// Otherwise, it defaults to "txt".
func DetectFileExtension(filename string) string {
	// Return ext without dot
	ext := FileExt(filename)

	// Return the extension if it is a common MIME type.
	if _, ok := CommonMimeTypes[ext]; ok {
		return ext
	}

	return "txt"

}

// DetectFileExtension returns the file extension for a given filename.
// If the extension is a common MIME type, it returns the extension.
// Otherwise, it defaults to "txt".
func FileExt(filename string) string {
	// Return ext without dot
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")

}
