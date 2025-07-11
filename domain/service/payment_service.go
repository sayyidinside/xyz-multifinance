package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/redis"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type PaymentService interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
	Create(ctx context.Context, input *model.PaymentInput) helpers.BaseResponse
}

type paymentService struct {
	paymentRepository     repository.PaymentRepository
	userRepository        repository.UserRepository
	transactionRepository repository.TransactionRepository
	installmentRepository repository.InstallmentRepository
	limitRepository       repository.LimitRepository
	lockRedis             *redis.LockClient
}

func NewPaymentService(
	paymentRepository repository.PaymentRepository,
	userRepository repository.UserRepository,
	transactionRepository repository.TransactionRepository,
	installmentRepository repository.InstallmentRepository,
	limitRepository repository.LimitRepository,
	lockRedis *redis.LockClient,
) PaymentService {
	return &paymentService{
		paymentRepository:     paymentRepository,
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		installmentRepository: installmentRepository,
		limitRepository:       limitRepository,
		lockRedis:             lockRedis,
	}
}

func (s *paymentService) GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	payment, err := s.paymentRepository.FindByUUID(ctx, uuid)
	if err != nil || payment == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	if !helpers.SelfOrAdminOnly(ctx, payment.Transaction.UserID) {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Message: "Unauthorized to access this data",
		})
	}

	paymentModel := model.PaymentToDetailModel(payment)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Payment data found",
		Data:    paymentModel,
	})
}

func (s *paymentService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var payments *[]entity.Payment
	var totalData int64
	var err error

	is_admin := ctx.Value(helpers.CtxKeyIsAdmin).(bool)
	if is_admin {
		payments, err = s.paymentRepository.FindAll(ctx, query)
		totalData = s.paymentRepository.Count(ctx, query)
	} else {
		session_user_id := ctx.Value(helpers.CtxKeyUserID).(float64)
		payments, err = s.paymentRepository.FindAllByUserID(ctx, query, uint(session_user_id))
		totalData = s.paymentRepository.CountByUserID(ctx, query, uint(session_user_id))
	}
	if err != nil || payments == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	paymentModels := model.PaymentToListModels(*payments)

	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Payment data found",
		Data:    paymentModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}

func (s *paymentService) Create(ctx context.Context, input *model.PaymentInput) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	paymentEntity, err := input.ToEntity()
	if err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Failed parsing input",
		})
	}

	session_user_id, ok := ctx.Value(helpers.CtxKeyUserID).(float64)
	if session_user_id == 0 || !ok {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Missing user id",
		})
	}

	user, err := s.userRepository.FindByID(ctx, uint(session_user_id))
	if err != nil || user == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "User Not Found",
			Errors:  err,
		})
	}

	installment, err := s.installmentRepository.FindByID(ctx, paymentEntity.InstallmentID)
	if err != nil || installment == nil || installment.PaymentStatus == entity.PaymentStatusFailed {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Installment Not Found",
			Errors:  err,
		})
	}

	transaction, err := s.transactionRepository.FindByID(ctx, installment.TransactionID)
	if err != nil || transaction == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	paymentEntity.TransactionID = transaction.ID

	// Lock installment
	installment_lock_name := fmt.Sprintf("lock:installment:%s", installment.UUID)
	lock_ttl := 10 * time.Second
	acquireinstallment, err := s.lockRedis.AcquireLock(ctx, installment_lock_name, lock_ttl)
	if !acquireinstallment || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "failed to acquire lock",
			Errors:  err,
		})
	}
	defer s.lockRedis.ReleaseLock(ctx, installment_lock_name)

	tx := s.paymentRepository.BeginTransaction(ctx)
	defer tx.Rollback()

	if err := s.paymentRepository.InsertWithTransaction(ctx, tx, paymentEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating transaction data",
			Errors:  logData.Err,
		})
	}

	newAmountPaid := installment.AmountPaid.Add(paymentEntity.Amount)
	if newAmountPaid.GreaterThanOrEqual(installment.AmountDue) {
		installment.AmountPaid = installment.AmountDue
		installment.PaymentStatus = entity.PaymentStatusPaid
		installment.PaidAt = sql.NullTime{Time: time.Now(), Valid: true}
	} else {
		installment.AmountPaid = newAmountPaid
		installment.PaymentStatus = entity.PaymentStatusPartial
	}

	if err := s.installmentRepository.UpdateWithTransaction(ctx, tx, installment); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating installments data",
			Errors:  logData.Err,
		})
	}

	if transaction.Tenor == installment.InstallmentNumber && installment.PaymentStatus == entity.PaymentStatusPaid {
		// Lock limit
		user_limit_lock_name := fmt.Sprintf("lock:userLimit:%s", user.UUID)
		lock_ttl := 10 * time.Second
		acquireUserLimit, err := s.lockRedis.AcquireLock(ctx, user_limit_lock_name, lock_ttl)
		if !acquireUserLimit || err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "failed to acquire lock",
				Errors:  err,
			})
		}
		defer s.lockRedis.ReleaseLock(ctx, user_limit_lock_name)

		transaction.Status = entity.TransactionPaid
		if err := s.transactionRepository.UpdateWithTransaction(ctx, tx, transaction); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating installments data",
				Errors:  logData.Err,
			})
		}

		limits, err := s.limitRepository.FindAllByUserID(ctx, &model.QueryGet{Limit: "100"}, user.ID)
		if err != nil || limits == nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "User limit not found",
				Errors:  logData.Err,
			})
		}
		updatedLimits := make([]entity.Limit, len(*limits))

		for i, limit := range *limits {
			newLimit := limit.CurrentLimit.Add(transaction.OnTheRoad)
			if newLimit.GreaterThan(limit.OriginalLimit) {
				newLimit = limit.OriginalLimit
			}

			updatedLimits[i] = entity.Limit{
				ID:            limit.ID,
				UserID:        limit.UserID,
				OriginalLimit: limit.OriginalLimit,
				CurrentLimit:  newLimit,
				Tenor:         limit.Tenor,
				UpdatedAt:     time.Now(),
			}
		}

		if err := s.limitRepository.BulkUpdateWithTransaction(ctx, tx, updatedLimits); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating installments data",
				Errors:  logData.Err,
			})
		}
	}

	tx.Commit()
	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: "Payment succesffully created",
	})
}
