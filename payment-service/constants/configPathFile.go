package constants

import (
	"fmt"
	"os"
)

const (
	FieldImagePathPrefix = "/api/v1/invoices/"
)

func BuildInvoiceURL(fileName string) string {
	baseURL := os.Getenv("BASE_URL") // ex: http://localhost:8080
	if baseURL == "" {
		baseURL = "http://localhost:8003" // fallback
	}
	return fmt.Sprintf("%s%s%s", baseURL, FieldImagePathPrefix, fileName)
}
