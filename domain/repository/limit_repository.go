package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type LimitRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	FindByID(ctx context.Context, id uint) (*entity.Limit, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Limit, error)
	FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Limit, error)
	FindAllByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (*[]entity.Limit, error)
	Count(ctx context.Context, query *model.QueryGet) int64
	CountUnscoped(ctx context.Context, query *model.QueryGet) int64
	Insert(ctx context.Context, limit *entity.Limit) error
	BulkInsertWithTransaction(ctx context.Context, tx *gorm.DB, limits []entity.Limit) error
	Update(ctx context.Context, limit *entity.Limit) error
	Delete(ctx context.Context, limit *entity.Limit) error
}

type limitRepository struct {
	*gorm.DB
}

func NewLimitRepository(db *gorm.DB) LimitRepository {
	return &limitRepository{DB: db}
}

func (r *limitRepository) FindByID(ctx context.Context, id uint) (*entity.Limit, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var limit entity.Limit
	result := r.DB.WithContext(ctx).
		Limit(1).
		Where("id = ?", id).
		Find(&limit)

	if result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &limit, nil
}

func (r *limitRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Limit, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var limit entity.Limit
	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Find(&limit); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return &limit, nil
}

func (r *limitRepository) FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Limit, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var limits []entity.Limit
	tx := r.DB.WithContext(ctx).Model(&entity.Limit{})

	var allowedFields = map[string]string{
		"tenor":   "limits.tenor",
		"created": "limits.created_at",
		"updated": "limits.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&limits).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &limits, nil
}

func (r *limitRepository) FindAllByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (*[]entity.Limit, error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var limits []entity.Limit
	tx := r.DB.WithContext(ctx).Model(&entity.Limit{}).
		Where("user_id = ?", user_id)

	var allowedFields = map[string]string{
		"tenor":   "limits.tenor",
		"created": "limits.created_at",
		"updated": "limits.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&limits).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return &limits, nil
}

func (r *limitRepository) Count(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Limit{})

	var allowedFields = map[string]string{
		"tenor":   "limits.tenor",
		"created": "limits.created_at",
		"updated": "limits.updated_at",
	}

	tx = tx.Scopes(
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return total
}

func (r *limitRepository) CountUnscoped(ctx context.Context, query *model.QueryGet) int64 {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.Limit{})

	var allowedFields = map[string]string{
		"tenor":   "limits.tenor",
		"created": "limits.created_at",
		"updated": "limits.updated_at",
	}

	tx = tx.Scopes(
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Count(&total).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
	}

	return total
}

func (r *limitRepository) Insert(ctx context.Context, limit *entity.Limit) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Create(limit).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}

func (r *limitRepository) Update(ctx context.Context, limit *entity.Limit) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Where("id = ?", limit.ID).Updates(limit).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}

func (r *limitRepository) Delete(ctx context.Context, limit *entity.Limit) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := r.DB.WithContext(ctx).Delete(limit).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}

	return nil
}

func (r *limitRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.DB.Begin()
}

func (r *limitRepository) BulkInsertWithTransaction(ctx context.Context, tx *gorm.DB, limits []entity.Limit) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := tx.WithContext(ctx).Create(limits).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return err
	}
	return nil
}
