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

const maxUploadSize = 10 << 20

var allowedUploadTypes = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
}

func UploadFile(c *gin.Context, formName string) (string, error) {

	file, err := c.FormFile(formName)

	if err != nil {
		return "", err
	}
	if file.Size <= 0 || file.Size > maxUploadSize {
		return "", fmt.Errorf("ukuran lampiran harus antara 1 byte dan %d MB", maxUploadSize/(1<<20))
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "." || strings.ContainsAny(file.Filename, `/\\`) {
		return "", errors.New("nama file lampiran tidak valid")
	}
	opened, err := file.Open()
	if err != nil {
		return "", err
	}
	defer opened.Close()
	header := make([]byte, 512)
	read, err := opened.Read(header)
	if err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(header[:read])
	allowedExt, ok := allowedUploadTypes[mimeType]
	if !ok || allowedExt != ext {
		return "", errors.New("tipe file lampiran tidak diizinkan")
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	if err := os.MkdirAll(uploadDir, 0o750); err != nil {
		return "", err
	}

	filename := uuid.New().String() + allowedExt

	path := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	return filename, nil
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
