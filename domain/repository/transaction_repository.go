package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	BeginTransaction(ctx context.Context) *gorm.DB
	FindByUUID(ctx context.Context, uuid uuid.UUID) (transaction *entity.Transaction, err error)
	FindAll(ctx context.Context, query *model.QueryGet) (transactions *[]entity.Transaction, err error)
	FindAllByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (transactions *[]entity.Transaction, err error)
	Count(ctx context.Context, query *model.QueryGet) (total int64)
	CountByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (total int64)
	CountUnscoped(ctx context.Context, query *model.QueryGet) (total int64)
	InsertWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Transaction) error
	UpdateWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Transaction) error
	DeleteWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Transaction) error
}

type transactionRepository struct {
	*gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{DB: db}
}

func (r *transactionRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.DB.Begin()
}

func (r *transactionRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (transaction *entity.Transaction, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Preload("User").Preload("User.Profile").Preload("Installments").Preload("Payments").
		Find(&transaction); result.Error != nil || result.RowsAffected == 0 {
		logData.Message = "Not Passed"
		logData.Err = result.Error
		return nil, result.Error
	}

	return
}

func (r *transactionRepository) FindAll(ctx context.Context, query *model.QueryGet) (transactions *[]entity.Transaction, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Transaction{}).
		Preload("User").Preload("User.Profile")

	var allowedFields = map[string]string{
		"created": "transactions.created_at",
		"updated": "transactions.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&transactions).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return
}

func (r *transactionRepository) FindAllByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (transactions *[]entity.Transaction, err error) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Transaction{}).Where("user_id = ?", user_id).
		Preload("User").Preload("User.Profile")

	var allowedFields = map[string]string{
		"created": "transactions.created_at",
		"updated": "transactions.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&transactions).Error; err != nil {
		logData.Message = "Not Passed"
		logData.Err = err
		return nil, err
	}

	return
}

func (r *transactionRepository) Count(ctx context.Context, query *model.QueryGet) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Transaction{})

	var allowedFields = map[string]string{
		"created": "transactions.created_at",
		"updated": "transactions.updated_at",
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

	return
}

func (r *transactionRepository) CountByUserID(ctx context.Context, query *model.QueryGet, user_id uint) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Transaction{}).Where("user_id = ?", user_id)

	var allowedFields = map[string]string{
		"created": "transactions.created_at",
		"updated": "transactions.updated_at",
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

	return
}

func (r *transactionRepository) CountUnscoped(ctx context.Context, query *model.QueryGet) (total int64) {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	tx := r.DB.WithContext(ctx).Model(&entity.Transaction{}).Unscoped()

	var allowedFields = map[string]string{
		"created": "transactions.created_at",
		"updated": "transactions.updated_at",
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

	return
}

func (r *transactionRepository) InsertWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Transaction) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	return tx.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) UpdateWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Transaction) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := tx.WithContext(ctx).Where("id = ?", transaction.ID).Updates(transaction).
		Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}

func (r *transactionRepository) DeleteWithTransaction(ctx context.Context, tx *gorm.DB, transaction *entity.Transaction) error {
	logData := helpers.CreateLog(r)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if err := tx.WithContext(ctx).Where("id = ?", transaction.ID).Delete(transaction).Error; err != nil {
		logData.Err = err
		logData.Message = "Not Passed"
		return err
	}

	return nil
}
