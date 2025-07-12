package handler

import (
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type DocumentHandler interface {
	UpdateDocument(c *fiber.Ctx) error
	GetDocument(c *fiber.Ctx) error
}

type documentHandler struct {
	service service.DocumentService
}

func NewDocumentHandler(service service.DocumentService) DocumentHandler {
	return &documentHandler{
		service: service,
	}
}

func (h *documentHandler) GetDocument(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	response := h.service.GetDocument(ctx)
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *documentHandler) UpdateDocument(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse

	ktpFile, err := c.FormFile("ktp_file")
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Missing ktp",
			Log:     &logData,
			Errors:  err,
		})
	}

	selfieFile, err := c.FormFile("selfie_file")
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Missing selfie",
			Log:     &logData,
			Errors:  err,
		})
	}

	ktpPath, err := h.saveUploadedFile(c, ktpFile)
	if err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Failed to save KTP file",
			Errors:  err,
		})
	}

	selfiePath, err := h.saveUploadedFile(c, selfieFile)
	if err != nil {
		os.Remove(ktpPath)
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Failed to save selfie file",
			Errors:  err,
		})
	}

	response = h.service.UpdateDocument(ctx, ktpPath, selfiePath)
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *documentHandler) saveUploadedFile(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".bin"
	}
	filename := uuid.New().String() + ext
	filePath := filepath.Join("storage", "uploads", filename)
	folderPath := filepath.Join("storage", "uploads")

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", err
	}

	if err := c.SaveFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
