package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"
)

func uploadProcces(image multipart.FileHeader) error {
	file, err := image.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil

}

func UploadImageLocal(images []multipart.FileHeader) ([]string, error) {

	photoPaths := make([]string, 0, len(images))
	for _, image := range images {
		err := uploadProcces(image)
		if err != nil {
			return nil, fmt.Errorf("failed to looping image: %w", err)
		}

		pathFile := fmt.Sprintf("images/%s-%s-%s", time.Now().Format("20060102150405"), image.Filename, path.Ext(image.Filename))

		basePath, err := filepath.Abs("/assets/")
		if err != nil {
			return nil, fmt.Errorf("gagal memendapatkan path absolut: %w", err)
		}

		photoPath := filepath.Join(basePath, "images", "field-images", pathFile)

		photoPaths = append(photoPaths, photoPath)

		out, err := os.Create(pathFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		log.Println("File saved successfully:", pathFile)

		defer out.Close()
	}

	return photoPaths, nil
}
