package utils

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context, formName string) (string, error) {

	file, err := c.FormFile(formName)

	if err != nil {
		return "", err
	}

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	ext := filepath.Ext(file.Filename)

	filename := uuid.New().String() + ext

	path := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	return filename, nil
}
