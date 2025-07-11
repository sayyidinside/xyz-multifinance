package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type TransactionHandler interface {
	Create(c *fiber.Ctx) error
	Cancel(c *fiber.Ctx) error
	GetTransaction(c *fiber.Ctx) error
	GetAllTransaction(c *fiber.Ctx) error
}

type transactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) TransactionHandler {
	return &transactionHandler{
		service: service,
	}
}

func (h *transactionHandler) Create(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var response helpers.BaseResponse
	var input model.TransactionInput

	if err := c.BodyParser(&input); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		input.Sanitize()
		if err := helpers.ValidateInput(input); err != nil {
			response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Log:     &logData,
				Errors:  err,
			})
		} else {
			response = h.service.Create(ctx, &input)
			response.Log = &logData
		}
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *transactionHandler) Cancel(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var response helpers.BaseResponse

	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid UUID Format",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		response = h.service.DeleteByUUID(ctx, uuid)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *transactionHandler) GetTransaction(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var response helpers.BaseResponse
	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID Format",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		response = h.service.GetByUUID(ctx, uuid)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *transactionHandler) GetAllTransaction(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var response helpers.BaseResponse
	query := new(model.QueryGet)

	if err := c.QueryParser(query); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request query",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		model.SanitizeQueryGet(query)

		url := c.BaseURL() + c.OriginalURL()
		response = h.service.GetAll(ctx, query, url)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}
