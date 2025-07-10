package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.UserProfile, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.UserProfile, error)
	FindByUserID(ctx context.Context, user_id uint) (*entity.UserProfile, error)
	Insert(ctx context.Context, profile *entity.UserProfile) error
	Update(ctx context.Context, profile *entity.UserProfile) error
	NikExist(ctx context.Context, profile *entity.UserProfile) bool
}

type profileRepository struct {
	*gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{DB: db}
}

func (r *profileRepository) FindByID(ctx context.Context, id uint) (*entity.UserProfile, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var profile entity.UserProfile
	result := r.DB.WithContext(ctx).
		Limit(1).
		Where("id = ?", id).
		Find(&profile)

	if result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &profile, nil
}

func (r *profileRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.UserProfile, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var profile entity.UserProfile
	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Find(&profile); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &profile, nil
}

func (r *profileRepository) FindByUserID(ctx context.Context, user_id uint) (*entity.UserProfile, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var profile entity.UserProfile
	if result := r.DB.WithContext(ctx).Limit(1).Where("user_id = ?", user_id).
		Find(&profile); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &profile, nil
}

func (r *profileRepository) Insert(ctx context.Context, profile *entity.UserProfile) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Create(profile).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}

func (r *profileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", profile.ID).Updates(profile).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}

func (r *profileRepository) Delete(ctx context.Context, profile *entity.UserProfile) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Delete(profile).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *profileRepository) NikExist(ctx context.Context, profile *entity.UserProfile) bool {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var totalData int64

	tx := r.DB.WithContext(ctx).Model(&entity.UserProfile{}).Where("nik = ?", profile.Nik)
	if profile.ID != 0 {
		tx = tx.Not("id = ?", profile.ID)
	}

	if err := tx.Count(&totalData).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
	}
	return totalData != 0
}
