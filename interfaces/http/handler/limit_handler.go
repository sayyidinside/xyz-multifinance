package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type LimitHandler interface {
	GetLimit(c *fiber.Ctx) error
}

type limitHandler struct {
	service service.LimitService
}

func NewLimitHandler(service service.LimitService) LimitHandler {
	return &limitHandler{
		service: service,
	}
}

func (h *limitHandler) GetLimit(c *fiber.Ctx) error {
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

		return helpers.ResponseFormatter(c, response)
	}

	uuid, err := uuid.Parse(c.Params("user_uuid"))
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID Format",
			Log:     &logData,
			Errors:  err,
		})

		return helpers.ResponseFormatter(c, response)
	}

	model.SanitizeQueryGet(query)

	url := c.BaseURL() + c.OriginalURL()
	response = h.service.GetUserLimit(ctx, uuid, query, url)
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}
