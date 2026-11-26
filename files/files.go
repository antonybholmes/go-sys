package files

import (
	"archive/zip"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/antonybholmes/go-sys/log"
)

var (
	// CommonMimeTypes = map[string]string{
	// 	".zip":  "application/zip",
	// 	".html": "text/html; charset=utf-8",
	// 	".json": "application/json",
	// 	".txt":  "text/plain; charset=utf-8",
	// 	".csv":  "text/csv; charset=utf-8",
	// 	".md":   "text/markdown; charset=utf-8",
	// 	".webp": "image/webp",
	// 	".svg":  "image/svg+xml",
	// 	".pdf":  "application/pdf",
	// }

	FileModeOwnerFullGroupOtherReadEnter os.FileMode = 0755
)

// DetectContentType determines the MIME type of a file based on its extension.
// Returns the MIME type as a string. If the extension is not recognized, it defaults to "application/octet-stream".
func DetectContentType(filename string) string {
	if mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(filename))); mimeType != "" {
		return mimeType
	}

	return "application/octet-stream"
}

// DetectFileExtension returns the file extension for a given filename.
// If the extension is a common MIME type, it returns the extension.
// Otherwise, it defaults to "txt".
func FileExt(filename string) string {
	// Return ext without dot
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")

}

// FileExtOr returns the file extension for a given filename.
// If the extension is empty, it returns the provided default
// extension. Otherwise, it returns the actual file extension.
func FileExtOr(filename, defaultExt string) string {
	ext := FileExt(filename)
	if ext == "" {
		return defaultExt
	}
	return ext
}

// Open a file for reading and return the file handle and its FileInfo.
func OpenFile(path string) (*os.File, os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open local file")
		return nil, nil, err
	}

	// need file size for S3 upload
	stat, err := file.Stat()
	if err != nil {
		log.Error().Err(err).Msg("Failed to stat local file")
		return nil, nil, err
	}

	return file, stat, nil
}

// AddFileToZip adds a file to the given zip archive. The file is added with its base name.
func AddFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer file.Close()

	return AddReaderToZip(zipWriter, filePath, file)
}

// Add contents of a reader to the zip file with the given file path as the entry name.
func AddReaderToZip(zipWriter *zip.Writer, filePath string, reader io.Reader) error {

	writer, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("create zip entry %s: %w", filePath, err)
	}

	if _, err := io.Copy(writer, reader); err != nil {
		return fmt.Errorf("copy %s: %w", filePath, err)
	}

	return nil
}
