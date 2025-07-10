package service

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type RegistrationService interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
	Create(ctx context.Context, input *model.UserRegisterInput) helpers.BaseResponse
	Activate(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse
}

type registrationService struct {
	repository     repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewRegistrationService(repository repository.UserRepository, roleRepository repository.RoleRepository) RegistrationService {
	return &registrationService{
		repository:     repository,
		roleRepository: roleRepository,
	}
}

func (s *registrationService) GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	user, err := s.repository.FindByUUID(ctx, uuid)
	if user == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User Not Found",
			Errors:  err,
		})
	}

	userModel := model.RegistrationToDetailModel(user)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User data found",
		Data:    userModel,
	})
}

func (s *registrationService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	query.FilterBy = "role"
	query.Filter = "User"

	users, err := s.repository.FindAll(ctx, query)
	if users == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User Not Found",
			Errors:  err,
		})
	}

	userModels := model.RegistrationToListModels(users)

	totalData := s.repository.Count(ctx, query)
	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User data found",
		Data:    userModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	})
}

func (s *registrationService) Create(ctx context.Context, input *model.UserRegisterInput) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	userEntity := input.ToEntity()

	if userEntity == nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
		})
	}

	// TODO: Use role id, make it more dynamicc
	userEntity.RoleID = 2

	if err := s.ValidateEntityInput(ctx, userEntity, input.RePassword); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	if err := s.repository.Insert(ctx, userEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating data",
			Errors:  logData.Err,
		})
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: "User successfully created",
	})

}

func (s *registrationService) Activate(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	user, err := s.repository.FindByUUID(ctx, uuid)
	if user == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User not found",
			Errors:  err,
		})
	}

	if user.ValidatedAt.Valid {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "User already activated",
			Errors:  err,
		})
	}

	user.ValidatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	user.UpdatedAt = time.Now()

	if err := s.repository.Update(ctx, user); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating data",
			Errors:  err,
		})
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User successfully activated",
	})
}

func (s *registrationService) ValidateEntityInput(ctx context.Context, user *entity.User, re_password string) interface{} {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	errors := []helpers.ValidationError{}
	errCh := make(chan helpers.ValidationError)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if role, err := s.roleRepository.FindByID(ctx, user.RoleID); role == nil || err != nil {
			errCh <- helpers.ValidationError{Field: "role_id", Tag: "not_found"}
			// errors = append(errors, helpers.ValidationError{
			// 	Field: "role_id",
			// 	Tag:   "not_found",
			// })
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if user.Password != re_password {
			errCh <- helpers.ValidationError{Field: "password", Tag: "not_match"}
			errCh <- helpers.ValidationError{Field: "re_password", Tag: "not_match"}
			// errors = append(errors, helpers.ValidationError{
			// 	Field: "password",
			// 	Tag:   "not_match",
			// })
			// errors = append(errors, helpers.ValidationError{
			// 	Field: "re_password",
			// 	Tag:   "not_match",
			// })
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if exist := s.repository.EmailExist(ctx, user); exist {
			errCh <- helpers.ValidationError{Field: "email", Tag: "duplicate"}
			// errors = append(errors, helpers.ValidationError{
			// 	Field: "email",
			// 	Tag:   "duplicate",
			// })
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if exist := s.repository.UsernameExist(ctx, user); exist {
			errCh <- helpers.ValidationError{Field: "username", Tag: "duplicate"}
			// errors = append(errors, helpers.ValidationError{
			// 	Field: "username",
			// 	Tag:   "duplicate",
			// })
		}
	}()

	// Close channel after all goroutines complete
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Collect errors from channel
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		logData.Message = "Validation error"
		logData.Err = errors
		return errors
	}
	return nil
}
