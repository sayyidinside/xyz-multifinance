package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type UserHandler interface {
	SuspendUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	GetAllUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) GetUser(c *fiber.Ctx) error {
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

func (h *userHandler) GetAllUser(c *fiber.Ctx) error {
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

func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)

	var input model.UserInput
	var response helpers.BaseResponse

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

func (h *userHandler) UpdateUser(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &logData,
			Errors:  err,
		})
	}

	var input model.UserUpdateInput
	if err := c.BodyParser(&input); err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		model.SanitizeUserUpdateInput(&input)

		if err := helpers.ValidateInput(input); err != nil {
			response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Log:     &logData,
				Errors:  err,
			})
		} else {
			response = h.service.UpdateByID(ctx, &input, uint(id))
			response.Log = &logData
		}

	}
	return helpers.ResponseFormatter(c, response)
}

func (h *userHandler) ResetPassword(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &logData,
			Errors:  err,
		})
	}

	var input model.ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		response = helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
			Errors:  err,
		}
	} else {
		model.SanitizeChangePasswordInput(&input)

		if err := helpers.ValidateInput(input); err != nil {
			response = helpers.BaseResponse{
				Status:  fiber.StatusBadRequest,
				Success: false,
				Message: "Invalid or malformed request body",
				Log:     &logData,
				Errors:  err,
			}
		} else {
			response = h.service.ChangePassByID(ctx, &input, uint(id))
			response.Log = &logData
		}
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *userHandler) DeleteUser(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		response = helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid ID format",
			Log:     &logData,
			Errors:  err,
		})
	} else {
		response = h.service.DeleteByID(ctx, uint(id))
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}

func (h *userHandler) SuspendUser(c *fiber.Ctx) error {
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
		response = h.service.SuspendByUUID(ctx, uuid)
		response.Log = &logData
	}

	return helpers.ResponseFormatter(c, response)
}
