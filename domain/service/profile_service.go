package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type ProfileService interface {
	UpdateProfile(ctx context.Context, input *model.UserProfileUpdate, id uint) helpers.BaseResponse
}

type profileService struct {
	profileRepository repository.ProfileRepository
	userRepository    repository.UserRepository
}

func NewProfileService(userRepository repository.UserRepository, profileRepository repository.ProfileRepository) ProfileService {
	return &profileService{
		profileRepository: profileRepository,
		userRepository:    userRepository,
	}
}

func (s *profileService) UpdateProfile(ctx context.Context, input *model.UserProfileUpdate, user_id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	profile, err := s.profileRepository.FindByUserID(ctx, user_id)
	if err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User not found",
			Errors:  err,
		})
	}

	profileEntity, err := input.ToEntity()
	if profileEntity == nil || err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
			Errors:  err,
		})
	}

	profileEntity.UserID = user_id
	if profile != nil {
		profileEntity.ID = profile.ID
	}

	if err := s.ValidateEntityInput(ctx, profileEntity); err != nil {
		return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		})
	}

	if profile == nil {
		if err := s.profileRepository.Insert(ctx, profileEntity); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating data",
				Errors:  err,
			})
		}
	} else {
		if err := s.profileRepository.Update(ctx, profileEntity); err != nil {
			return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
				Status:  fiber.StatusInternalServerError,
				Success: false,
				Message: "Error updating data",
				Errors:  err,
			})
		}
	}

	return helpers.LogBaseResponse(&logData, helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Profile successfully updated",
	})
}

func (s *profileService) ValidateEntityInput(ctx context.Context, user *entity.UserProfile) interface{} {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	errors := []helpers.ValidationError{}

	if exist := s.profileRepository.NikExist(ctx, user); exist {
		errors = append(errors, helpers.ValidationError{
			Field: "nik",
			Tag:   "duplicate",
		})
	}

	if len(errors) > 0 {
		logData.Message = "Validation error"
		logData.Err = errors
		return errors
	}
	return nil
}
