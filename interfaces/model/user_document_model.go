package model

import "github.com/sayyidinside/gofiber-clean-fresh/domain/entity"

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
		SelfieFile: userDocument.SelfieFile,
		KtpFile:    userDocument.KtpFile,
	}
}
