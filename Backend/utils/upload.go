package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const MaxUploadSize = 10 << 20

var allowedUploadTypes = map[string]map[string]bool{
	"application/pdf": {".pdf": true},
	"image/jpeg":      {".jpg": true, ".jpeg": true},
	"image/png":       {".png": true},
}

type UploadedFile struct {
	Path     string
	Original string
	MIMEType string
	Size     int64
}

func UploadFileMetadata(c *gin.Context, formName string) (UploadedFile, error) {
	var result UploadedFile

	file, err := c.FormFile(formName)

	if err != nil {
		return result, err
	}
	if file.Size <= 0 || file.Size > MaxUploadSize {
		return result, fmt.Errorf("ukuran lampiran harus antara 1 byte dan %d MB", MaxUploadSize/(1<<20))
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "." || strings.ContainsAny(file.Filename, `/\\`) {
		return result, errors.New("nama file lampiran tidak valid")
	}
	opened, err := file.Open()
	if err != nil {
		return result, err
	}
	defer opened.Close()
	header := make([]byte, 512)
	read, err := opened.Read(header)
	if err != nil && read == 0 {
		return result, err
	}
	mimeType := http.DetectContentType(header[:read])
	allowedExts, ok := allowedUploadTypes[mimeType]
	if !ok || !allowedExts[ext] {
		return result, errors.New("tipe file lampiran tidak diizinkan")
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	if err := os.MkdirAll(uploadDir, 0o750); err != nil {
		return result, err
	}

	filename := uuid.New().String() + ext

	path := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		return result, err
	}

	return UploadedFile{Path: filename, Original: filepath.Base(file.Filename), MIMEType: mimeType, Size: file.Size}, nil
}

func UploadFile(c *gin.Context, formName string) (string, error) {
	file, err := UploadFileMetadata(c, formName)
	return file.Path, err
}

func DeleteUploadedFile(filename string) error {
	if filename == "" || filepath.Base(filename) != filename {
		return errors.New("nama file lampiran tidak valid")
	}
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	err := os.Remove(filepath.Join(uploadDir, filename))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
