package model

import (
	"path/filepath"
	"strings"

	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
)

type (
	UserDocumentDetail struct {
		ID         uint   `json:"id"`
		UserID     uint   `json:"user_id"`
		SelfieFile string `json:"selfie_file"`
		KtpFile    string `json:"ktp_file"`
	}
)

func UserDocumentToDetailModel(userDocument *entity.UserDocument) *UserDocumentDetail {

	return &UserDocumentDetail{
		ID:         userDocument.ID,
		UserID:     userDocument.UserID,
		SelfieFile: filePathToURL(userDocument.SelfieFile),
		KtpFile:    filePathToURL(userDocument.KtpFile),
	}
}

func filePathToURL(filePath string) string {
	if filePath == "" {
		return ""
	}

	baseURL := config.AppConfig.BaseUrl
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}

	filename := filepath.Base(filePath)

	baseURL = strings.TrimSuffix(baseURL, "/")

	return baseURL + "/files/" + filename
}
