package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func UploadImageLocal(images []multipart.FileHeader) ([]string, error) {
	photoPaths := make([]string, 0, len(images))

	for _, image := range images {
		basePath, err := filepath.Abs("assets/")
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan path absolut: %w", err)
		}

		// Pastikan folder penyimpanan ada
		folderPath := filepath.Join(basePath, "images", "field-images")
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

		// Tentukan path file
		photoPath := filepath.Join(folderPath, fmt.Sprintf("%s-%s%s",
			time.Now().Format("20060102150405"), image.Filename, filepath.Ext(image.Filename)))

		// Buat file
		out, err := os.Create(photoPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		defer out.Close()

		// Buka file dari multipart
		file, err := image.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open image: %w", err)
		}
		defer file.Close()

		// Salin isi file
		_, err = io.Copy(out, file)
		if err != nil {
			return nil, fmt.Errorf("failed to save image: %w", err)
		}

		photoPaths = append(photoPaths, photoPath)
	}

	return photoPaths, nil
}
