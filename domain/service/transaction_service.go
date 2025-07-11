package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/redis"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"github.com/shopspring/decimal"
)

type TransactionService interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
	GetAllByUserUUID(ctx context.Context, query *model.QueryGet, url string, uuid uuid.UUID) helpers.BaseResponse
	Create(ctx context.Context, input *model.TransactionInput) helpers.BaseResponse
	UpdateByUUID(ctx context.Context, input *model.TransactionInput, uuid uuid.UUID) helpers.BaseResponse
	DeleteByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
	userRepository        repository.UserRepository
	limitRepository       repository.LimitRepository
	installmentRepository repository.InstallmentRepository
	lockRedis             *redis.LockClient
}

func NewTransactionService(
	transactionRepository repository.TransactionRepository,
	userRepository repository.UserRepository,
	limitRepository repository.LimitRepository,
	installmentRepository repository.InstallmentRepository,
	lockRedis *redis.LockClient,
) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
		limitRepository:       limitRepository,
		installmentRepository: installmentRepository,
		lockRedis:             lockRedis,
	}
}

func (s *transactionService) GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	transaction, err := s.transactionRepository.FindByUUID(ctx, uuid)
	if err != nil || transaction == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	if !helpers.SelfOrAdminOnly(ctx, transaction.UserID) {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Message: "Unauthorized to access this data",
		})
	}

	transactionModel := model.TransactionToDetailModel(transaction)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Transaction data found",
		Data:    transactionModel,
	})
}

func (s *transactionService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var transactions *[]entity.Transaction
	var totalData int64
	var err error

	is_admin := ctx.Value(helpers.CtxKeyIsAdmin).(bool)
	if is_admin {
		transactions, err = s.transactionRepository.FindAll(ctx, query)
		totalData = s.transactionRepository.Count(ctx, query)
	} else {
		session_user_id := ctx.Value(helpers.CtxKeyUserID).(float64)
		transactions, err = s.transactionRepository.FindAllByUserID(ctx, query, uint(session_user_id))
		totalData = s.transactionRepository.CountByUserID(ctx, query, uint(session_user_id))
	}
	if err != nil || transactions == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	transactionModels := model.TransactionToListModels(*transactions)

	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Transactions data found",
		Data:    transactionModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}

func (s *transactionService) GetAllByUserUUID(ctx context.Context, query *model.QueryGet, url string, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	user, err := s.userRepository.FindByUUID(ctx, uuid)
	if user == nil || err != nil {
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

	transactions, err := s.transactionRepository.FindAllByUserID(ctx, query, user.ID)
	if err != nil || transactions == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	transactionModels := model.TransactionToListModels(*transactions)

	totalData := s.transactionRepository.CountByUserID(ctx, query, user.ID)
	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Transactions data found",
		Data:    transactionModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}

func (s *transactionService) Create(ctx context.Context, input *model.TransactionInput) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	transactionEntity, err := input.ToEntity()
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

	transactionEntity.UserID = user.ID

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

	// Limit calculation
	updatedLimits, errResponse := s.generateUpdatedLimitList(ctx, user.ID, transactionEntity.Tenor, transactionEntity.OnTheRoad, true)
	if errResponse != nil {
		return helpers.LogBaseResponse(&logData, *errResponse)
	}

	tx := s.transactionRepository.BeginTransaction(ctx)
	defer tx.Rollback()

	if err := s.transactionRepository.InsertWithTransaction(ctx, tx, transactionEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating transaction data",
			Errors:  logData.Err,
		})
	}

	// Generate installment
	newInstallments := make([]entity.TransactionInstallment, transactionEntity.Tenor)
	for month := range int(transactionEntity.Tenor) {
		dueDate := time.Now().AddDate(0, month+1, 0)
		newInstallments[month] = entity.TransactionInstallment{
			TransactionID:     transactionEntity.ID,
			InstallmentNumber: uint(month + 1),
			AmountDue:         transactionEntity.MonthlyInstallment,
			DueDate:           dueDate,
		}
	}

	if err := s.installmentRepository.BulkInsertWithTransaction(ctx, tx, newInstallments); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating installments data",
			Errors:  logData.Err,
		})
	}

	if err := s.limitRepository.BulkUpdateWithTransaction(ctx, tx, updatedLimits); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating limits data",
			Errors:  logData.Err,
		})
	}

	tx.Commit()
	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: "Transaction succesffully created",
	})
}

func (s *transactionService) UpdateByUUID(ctx context.Context, input *model.TransactionInput, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "To be changed",
	})
}

func (s *transactionService) DeleteByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	transaction, err := s.transactionRepository.FindByUUID(ctx, uuid)
	if err != nil || transaction == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Transaction Not Found",
			Errors:  err,
		})
	}

	if !helpers.SelfOrAdminOnly(ctx, transaction.UserID) {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusForbidden,
			Success: false,
			Message: "Unauthorized to access this data",
		})
	}

	if transaction.Status != entity.TransactionActive {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Transaction already paid or already cancelled",
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

	lock_ttl := 10 * time.Second

	// Lock limit
	user_limit_lock_name := fmt.Sprintf("lock:userLimit:%s", user.UUID)
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

	// Lock transaction
	transaction_lock_name := fmt.Sprintf("lock:transaction:%s", transaction.UUID)
	acquireTransaction, err := s.lockRedis.AcquireLock(ctx, transaction_lock_name, lock_ttl)
	if !acquireTransaction || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "failed to acquire lock",
			Errors:  err,
		})
	}
	defer s.lockRedis.ReleaseLock(ctx, transaction_lock_name)

	// Lock installment
	installment_lock_name := fmt.Sprintf("lock:installment:%s", transaction.UUID)
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

	// Limit calculation
	updatedLimits, errResponse := s.generateUpdatedLimitList(ctx, user.ID, transaction.Tenor, transaction.OnTheRoad, false)
	if errResponse != nil {
		return helpers.LogBaseResponse(&logData, *errResponse)
	}

	tx := s.transactionRepository.BeginTransaction(ctx)
	defer tx.Rollback()

	transaction.Status = entity.TransactionCanceled
	if err := s.transactionRepository.UpdateWithTransaction(ctx, tx, transaction); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating transaction data",
			Errors:  logData.Err,
		})
	}

	if err := s.installmentRepository.CancelWithTransaction(ctx, tx, transaction.ID); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating installments data",
			Errors:  logData.Err,
		})
	}

	if err := s.limitRepository.BulkUpdateWithTransaction(ctx, tx, updatedLimits); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating limits data",
			Errors:  logData.Err,
		})
	}

	tx.Commit()
	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Transaction successfully cancelled",
	})
}

func (s *transactionService) generateUpdatedLimitList(ctx context.Context, user_id uint, tenor uint, otr decimal.Decimal, is_reduce bool) ([]entity.Limit, *helpers.BaseResponse) {
	limits, err := s.limitRepository.FindAllByUserID(ctx, &model.QueryGet{}, user_id)
	if err != nil || limits == nil {
		return []entity.Limit{}, &helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "User limit not found",
			Errors:  err,
		}
	}

	tenor_limit_found := false
	updatedLimits := make([]entity.Limit, len(*limits))
	for i, limit := range *limits {
		if limit.Tenor == tenor {
			tenor_limit_found = true
			if limit.CurrentLimit.LessThan(otr) && is_reduce {
				return []entity.Limit{}, &helpers.BaseResponse{
					Status:  fiber.StatusBadRequest,
					Success: false,
					Message: "Overlimit",
					Errors:  err,
				}
			}
		}

		var newLimit decimal.Decimal
		if is_reduce {
			if newLimit.IsPositive() {
				newLimit = limit.CurrentLimit.Sub(otr)
				if newLimit.IsNegative() {
					newLimit = decimal.Zero
				}
			}
		} else {
			newLimit = limit.CurrentLimit.Add(otr)
			if newLimit.GreaterThan(limit.OriginalLimit) {
				newLimit = limit.OriginalLimit
			}
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
	if !tenor_limit_found {
		return []entity.Limit{}, &helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "User limit tenor not found",
			Errors:  err,
		}
	}

	return updatedLimits, nil
}
