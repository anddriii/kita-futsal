package constants

import (
	"fmt"
	"os"
)

const (
	FieldImagePathPrefix = "/assets/images/field-images/"
)

func BuildFullImagePath(fileName string) string {
	baseURL := os.Getenv("BASE_URL") // ex: http://localhost:8080
	if baseURL == "" {
		baseURL = "http://localhost:8002" // fallback
	}
	return fmt.Sprintf("%s%s%s", baseURL, FieldImagePathPrefix, fileName)
}
