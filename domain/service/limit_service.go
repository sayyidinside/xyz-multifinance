package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type LimitService interface {
	GetUserLimit(ctx context.Context, uuid uuid.UUID, query *model.QueryGet, url string) helpers.BaseResponse
}

type limitService struct {
	userRepository  repository.UserRepository
	limitRepository repository.LimitRepository
}

func NewLimitService(userRepository repository.UserRepository, limitRepository repository.LimitRepository) LimitService {
	return &limitService{
		userRepository:  userRepository,
		limitRepository: limitRepository,
	}
}

func (s *limitService) GetUserLimit(ctx context.Context, uuid uuid.UUID, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	user, err := s.userRepository.FindByUUID(ctx, uuid)
	if err != nil || user == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User Not Found",
			Errors:  err,
		})
	}

	if !helpers.SelfOrAdminOnly(ctx, user.ID) {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Message: "Unauthorized to access this data",
		})
	}

	limits, err := s.limitRepository.FindAllByUserID(ctx, query, user.ID)
	if err != nil || limits == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User limit Not Found",
			Errors:  err,
		})
	}

	limitModel := model.LimitToListModels(*limits)

	totalData := s.limitRepository.CountByUserID(ctx, query, user.ID)
	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User limit data found",
		Data:    limitModel,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}
