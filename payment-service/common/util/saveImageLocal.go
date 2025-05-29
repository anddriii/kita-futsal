package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UploadImageLocal(images []multipart.FileHeader) ([]string, error) {
	photoNames := make([]string, 0, len(images))

	for _, image := range images {
		basePath, err := filepath.Abs("assets/images/field-images")
		if err != nil {
			return nil, fmt.Errorf("failed to get abs path: %w", err)
		}

		if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

		FilenameRemoveSpace := strings.ReplaceAll(image.Filename, " ", "-")

		fileName := fmt.Sprintf("%s-%s", time.Now().Format("20060102150405"), FilenameRemoveSpace)
		fullPath := filepath.Join(basePath, fileName)

		out, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		defer out.Close()

		in, err := image.Open()
		if err != nil {
			return nil, err
		}
		defer in.Close()

		if _, err := io.Copy(out, in); err != nil {
			return nil, err
		}

		photoNames = append(photoNames, fileName)
	}

	return photoNames, nil
}
