package model

import (
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/utils"
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
		SelfieFile: utils.FilePathToURL(userDocument.SelfieFile),
		KtpFile:    utils.FilePathToURL(userDocument.KtpFile),
	}
}
