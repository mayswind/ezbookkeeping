package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

var imageFileExtensionContentTypeMap = map[string]string{
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"webp": "image/webp",
}

// GetImageContentType returns the content type of specified image file extension or returns empty when the file extension is not image or not supported
func GetImageContentType(fileExtension string) string {
	contentType, exists := imageFileExtensionContentTypeMap[fileExtension]

	if !exists {
		return ""
	}

	return contentType
}

// ListFileNamesWithPrefixAndSuffix returns file name list which has specified prefix and suffix
func ListFileNamesWithPrefixAndSuffix(path string, prefix string, suffix string) []string {
	dir, err := os.Open(path)

	if err != nil {
		return nil
	}

	fileInfos, err := dir.Readdir(0)

	if err != nil {
		return nil
	}

	var fileNames []string

	for i := 0; i < len(fileInfos); i++ {
		fileInfo := fileInfos[i]

		if !fileInfo.IsDir() &&
			strings.HasPrefix(fileInfo.Name(), prefix) &&
			strings.HasSuffix(fileInfo.Name(), suffix) {
			fileNames = append(fileNames, fileInfo.Name())
		}
	}

	return fileNames
}

// IsExists returns whether specified file or directory path exits
func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// WriteFile would write file according to specified content
func WriteFile(path string, data []byte) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	n, err := file.Write(data)

	if err == nil && n < len(data) {
		return io.ErrShortWrite
	}

	return err
}

// GetFileNameWithoutExtension returns the file name without extension
func GetFileNameWithoutExtension(path string) string {
	fileName := filepath.Base(path)
	extension := filepath.Ext(fileName)

	if len(extension) < 1 {
		return fileName
	}

	return fileName[0 : len(fileName)-len(extension)]
}

// GetFileNameExtension returns the file extension without dot
func GetFileNameExtension(path string) string {
	extension := filepath.Ext(path)

	if len(extension) < 1 || extension[0] != '.' {
		return extension
	}

	return extension[1:]
}

// IdentReader returns the original io reader
func IdentReader(encoding string, input io.Reader) (io.Reader, error) {
	return input, nil
}
