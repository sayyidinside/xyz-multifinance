package service

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type DocumentService interface {
	UpdateDocument(ctx context.Context, ktp_file string, selfie_file string) helpers.BaseResponse
	GetDocument(ctx context.Context) helpers.BaseResponse
}

type documentService struct {
	documentRepository repository.DocumentRepository
	userRepository     repository.UserRepository
}

func NewDocumentService(userRepository repository.UserRepository, DocumentRepository repository.DocumentRepository) DocumentService {
	return &documentService{
		documentRepository: DocumentRepository,
		userRepository:     userRepository,
	}
}

func (s *documentService) UpdateDocument(ctx context.Context, ktp_file string, selfie_file string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	session_user_id, ok := ctx.Value(helpers.CtxKeyUserID).(float64)
	if session_user_id == 0 || !ok {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Missing user id",
		})
	}

	document, err := s.documentRepository.FindByUserID(ctx, uint(session_user_id))
	if err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User not found",
			Errors:  err,
		})
	}

	if document == nil {
		// Create new document
		newDoc := &entity.UserDocument{
			UserID:     uint(session_user_id),
			KtpFile:    ktp_file,
			SelfieFile: selfie_file,
		}
		if err := s.documentRepository.Insert(ctx, newDoc); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating data",
				Errors:  err,
			})
		}
	} else {
		go s.deleteOldFiles(document.KtpFile, document.SelfieFile)

		document.KtpFile = ktp_file
		document.SelfieFile = selfie_file
		if err := s.documentRepository.Update(ctx, document); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating data",
				Errors:  err,
			})
		}
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Document successfully updated",
	})
}

func (s *documentService) GetDocument(ctx context.Context) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	session_user_id, ok := ctx.Value(helpers.CtxKeyUserID).(float64)
	if session_user_id == 0 || !ok {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Missing user id",
		})
	}

	Document, err := s.documentRepository.FindByUserID(ctx, uint(session_user_id))
	if err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Document not found",
		})
	}

	DocumentModel := model.UserDocumentToDetailModel(Document)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Document data found",
		Data:    DocumentModel,
	})
}

func (s *documentService) deleteOldFiles(paths ...string) {
	for _, path := range paths {
		if path != "" {
			os.Remove(path)
		}
	}
}
