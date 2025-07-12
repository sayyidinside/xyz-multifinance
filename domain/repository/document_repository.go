package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.UserDocument, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.UserDocument, error)
	FindByUserID(ctx context.Context, user_id uint) (*entity.UserDocument, error)
	Insert(ctx context.Context, document *entity.UserDocument) error
	Update(ctx context.Context, document *entity.UserDocument) error
}

type documentRepository struct {
	*gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{DB: db}
}

func (r *documentRepository) FindByID(ctx context.Context, id uint) (*entity.UserDocument, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var document entity.UserDocument
	result := r.DB.WithContext(ctx).
		Limit(1).
		Where("id = ?", id).
		Find(&document)

	if result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &document, nil
}

func (r *documentRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.UserDocument, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var document entity.UserDocument
	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Find(&document); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &document, nil
}

func (r *documentRepository) FindByUserID(ctx context.Context, user_id uint) (*entity.UserDocument, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var document entity.UserDocument
	if result := r.DB.WithContext(ctx).Limit(1).Where("user_id = ?", user_id).
		Find(&document); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &document, nil
}

func (r *documentRepository) Insert(ctx context.Context, document *entity.UserDocument) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Create(document).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}

func (r *documentRepository) Update(ctx context.Context, document *entity.UserDocument) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", document.ID).Updates(document).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}

func (r *documentRepository) Delete(ctx context.Context, document *entity.UserDocument) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Delete(document).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}
