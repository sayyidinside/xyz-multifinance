package utils

import (
	"path/filepath"
	"strings"

	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
)

func FilePathToURL(filePath string) string {
	if filePath == "" {
		return ""
	}

	baseURL := config.AppConfig.BaseUrl
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	filename := filepath.Base(filePath)
	filename = strings.ReplaceAll(filename, "..", "")

	baseURL = strings.TrimSuffix(baseURL, "/")

	return baseURL + "/files/" + filename
}
