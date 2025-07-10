package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type ProfileHandler interface {
	UpdateProfile(c *fiber.Ctx) error
}

type profileHandler struct {
	service service.ProfileService
}

func NewProfileHandler(service service.ProfileService) ProfileHandler {
	return &profileHandler{
		service: service,
	}
}

func (h *profileHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := helpers.ExtractIdentifierAndUsername(c)
	logData := helpers.CreateLog(h)

	defer helpers.LogSystemWithDefer(ctx, &logData)
	var response helpers.BaseResponse

	user_id := c.Locals("user_id").(float64)

	var input model.UserProfileUpdate
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
			response = h.service.UpdateProfile(ctx, &input, uint(user_id))
			response.Log = &logData
		}

	}
	return helpers.ResponseFormatter(c, response)
}
