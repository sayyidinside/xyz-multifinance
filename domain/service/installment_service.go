package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type InstallmentService interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
}

type installmentService struct {
	installmentRepository repository.InstallmentRepository
	userRepository        repository.UserRepository
}

func NewInstallmentService(
	installmentRepository repository.InstallmentRepository,
	userRepository repository.UserRepository,
) InstallmentService {
	return &installmentService{
		installmentRepository: installmentRepository,
		userRepository:        userRepository,
	}
}

func (s *installmentService) GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	installment, err := s.installmentRepository.FindByUUID(ctx, uuid)
	if err != nil || installment == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	if !helpers.SelfOrAdminOnly(ctx, installment.Transaction.UserID) {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Message: "Unauthorized to access this data",
		})
	}

	installmentModel := model.TransactionInstallmentToDetailModel(installment)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Installment data found",
		Data:    installmentModel,
	})
}

func (s *installmentService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var installments *[]entity.TransactionInstallment
	var totalData int64
	var err error

	is_admin := ctx.Value(helpers.CtxKeyIsAdmin).(bool)
	if is_admin {
		installments, err = s.installmentRepository.FindAll(ctx, query)
		totalData = s.installmentRepository.Count(ctx, query)
	} else {
		session_user_id := ctx.Value(helpers.CtxKeyUserID).(float64)
		installments, err = s.installmentRepository.FindAllByUserID(ctx, query, uint(session_user_id))
		totalData = s.installmentRepository.CountByUserID(ctx, query, uint(session_user_id))
	}
	if err != nil || installments == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	installmentModels := model.TransactionInstallmentToListModels(*installments)

	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Installment data found",
		Data:    installmentModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}
